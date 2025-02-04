use super::{AsNode, BumpStr, BumpVec, Field, Fields, Node};

/// HeapDoc is a document representation stored in the heap.
#[derive(Debug, rkyv::Archive, rkyv::Serialize)]
#[archive(archived = "ArchivedDoc")]
pub struct HeapDoc<'alloc> {
    /// Root node of the document.
    pub root: HeapNode<'alloc>,
    /// Arbitrary flags used to persist document processing status.
    pub flags: u8,
}

/// HeapNode is a document node representation stored in the heap.
// The additional archive bounds are required to satisfy the compiler due to
// the recursive nature of this structure. For more explanation see:
// https://github.com/rkyv/rkyv/blob/master/examples/json/src/main.rs
#[derive(Debug, rkyv::Archive, rkyv::Serialize)]
#[archive(
    archived = "ArchivedNode",
    bound(
        serialize = "__S: rkyv::ser::ScratchSpace + rkyv::ser::Serializer + rkyv::ser::SharedSerializeRegistry"
    )
)]
pub enum HeapNode<'alloc> {
    Array(#[omit_bounds] BumpVec<'alloc, HeapNode<'alloc>>),
    Bool(bool),
    Bytes(BumpVec<'alloc, u8>),
    Float(f64),
    NegInt(i64),
    Null,
    Object(BumpVec<'alloc, HeapField<'alloc>>),
    PosInt(u64),
    String(BumpStr<'alloc>),
}

/// HeapField is a field representation stored in the heap.
#[derive(Debug, rkyv::Archive, rkyv::Serialize)]
#[archive(
    archived = "ArchivedField",
    bound(
        serialize = "__S: rkyv::ser::ScratchSpace + rkyv::ser::Serializer + rkyv::ser::SharedSerializeRegistry"
    )
)]
pub struct HeapField<'alloc> {
    pub property: BumpStr<'alloc>,
    #[omit_bounds]
    pub value: HeapNode<'alloc>,
}

impl<'alloc> HeapNode<'alloc> {
    // new_allocator builds a bumpalo::Bump allocator for use in building HeapNodes.
    // It's a trivial helper which can reduce type imports.
    pub fn new_allocator() -> bumpalo::Bump {
        bumpalo::Bump::new()
    }
}

impl<'alloc> AsNode for HeapNode<'alloc> {
    type Fields = [HeapField<'alloc>];

    fn as_node<'a>(&'a self) -> Node<'a, Self> {
        match self {
            HeapNode::Array(a) => Node::Array(a),
            HeapNode::Bool(b) => Node::Bool(*b),
            HeapNode::Bytes(b) => Node::Bytes(b),
            HeapNode::Float(n) => Node::Number(json::Number::Float(*n)),
            HeapNode::NegInt(n) => Node::Number(json::Number::Signed(*n)),
            HeapNode::Null => Node::Null,
            HeapNode::Object(o) => Node::Object(o.as_slice()),
            HeapNode::PosInt(n) => Node::Number(json::Number::Unsigned(*n)),
            HeapNode::String(s) => Node::String(s),
        }
    }
}

impl<'alloc> Fields<HeapNode<'alloc>> for [HeapField<'alloc>] {
    type Field<'a> = &'a HeapField<'alloc> where 'alloc: 'a;
    type Iter<'a> = std::slice::Iter<'a, HeapField<'alloc>> where 'alloc: 'a;

    fn get<'a>(&'a self, property: &str) -> Option<Self::Field<'a>> {
        match self.binary_search_by(|l| l.property.cmp(property)) {
            Ok(ind) => Some(&self[ind]),
            Err(_) => None,
        }
    }

    fn len(&self) -> usize {
        self.len()
    }

    fn iter<'a>(&'a self) -> Self::Iter<'a> {
        self.iter()
    }
}

impl<'a, 'alloc> Field<'a, HeapNode<'alloc>> for &'a HeapField<'alloc> {
    fn property(&self) -> &'a str {
        &self.property
    }
    fn value(&self) -> &'a HeapNode<'alloc> {
        &self.value
    }
}

impl<'alloc> BumpVec<'alloc, HeapField<'alloc>> {
    /// Insert or obtain a mutable reference to a child HeapNode with the given property.
    pub fn insert_property(
        &mut self,
        property: &str,
        alloc: &'alloc bumpalo::Bump,
    ) -> &mut HeapNode<'alloc> {
        let ind = match self.binary_search_by(|l| l.property.cmp(property)) {
            Ok(ind) => ind,
            Err(ind) => {
                self.insert(
                    ind,
                    HeapField {
                        property: BumpStr::from_str(property, alloc),
                        value: HeapNode::Null,
                    },
                    alloc,
                );
                ind
            }
        };
        &mut self[ind].value
    }

    // Remove the named property, returning its removed HeapField if found.
    pub fn remove_property(&mut self, property: &str) -> Option<HeapField<'alloc>> {
        match self.binary_search_by(|l| l.property.cmp(property)) {
            Ok(ind) => Some(self.remove(ind)),
            Err(_) => None,
        }
    }
}
