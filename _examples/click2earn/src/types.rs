use serde::Serialize;
use serde_with::skip_serializing_none;

#[skip_serializing_none]
#[derive(Serialize, Default)]
pub struct Param {
    pub int32: Option<i32>,
    pub int64: Option<i64>,
    pub float32: Option<f32>,
    pub float64: Option<f64>,
    pub string: Option<String>,
    pub time: Option<String>, //  rfc3339 encoding
    pub bool: Option<bool>,
    pub bytes: Option<String>, // base64 encoding
}

#[derive(Serialize)]
pub struct DBQuery {
    pub statement: String,
    pub params: Vec<Param>,
}

#[derive(Serialize)]
pub struct Tx {
    pub to: String,
    pub data: String,
    pub value: String,
}
