use std::mem;
use std::str;

use anyhow::{bail, Context, Result};
use serde_json::Value;

mod sdk;
mod types;

#[no_mangle]
pub extern "C" fn alloc(size: i32) -> *mut u8 {
    let mut buf: Vec<u8> = Vec::with_capacity(size as _);
    let ptr = buf.as_mut_ptr();
    mem::forget(buf);
    return ptr;
}

static mut MSGNUM: u64 = 0;

#[no_mangle]
pub extern "C" fn start(resource_id: i32) -> i32 {
    unsafe {
        MSGNUM += 1;
    }

    if unsafe { MSGNUM == 1 } {
        let query = types::DBQuery {
            statement: String::from(
                "CREATE TABLE Click2Earn (
                Addr varchar(255) NOT NULL PRIMARY KEY,
                TotalClicks Int,
                TotalTx Int
            );",
            ),
            params: vec![],
        };
        if sdk::set_db(&serde_json::to_string(&query).unwrap()).is_err() {
            sdk::log_error("fail to create sql table");
            return -1;
        };
    }

    let data_u8 = match sdk::get_data(resource_id) {
        Some(data) => data,
        _ => {
            sdk::log_error("fail to get data");
            return -1;
        }
    };

    sdk::log_info(format!("receive message: {}", str::from_utf8(&data_u8).unwrap()).as_str());

    // sent tx
    let tx_hash =  sdk::send_tx(
            &String::from("0x1b1EAD130bD9162ca05E743131e1837DDE5a3196"),
            &String::from("0"),
            &format!("40c10f19000000000000000000000000{}0000000000000000000000000000000000000000000000000de0b6b3a7640000", "97186a21fa8e7955c0f154f960d588c3aca44f14"),
        ).ok();

    if upsert_record(data_u8, tx_hash).is_err() {
        sdk::log_error("fail to insert record");
        return -1;
    }

    if unsafe { MSGNUM % 5 == 0 } {
        // query
        let prestate = String::from("SELECT * FROM Click2Earn;");
        let query = types::DBQuery {
            statement: prestate,
            params: vec![],
        };
        let str = serde_json::to_string(&query).unwrap();
        let val = match sdk::get_db(&str) {
            Some(val) => val,
            _ => return -1,
        };
        sdk::log_info(str::from_utf8(&val).unwrap());
    }

    return 0;
}

fn upsert_record(data_u8: Vec<u8>, tx_hash: Option<String>) -> Result<()> {
    let data: Value = serde_json::from_slice(data_u8.as_slice())?;

    let addr = data
        .get("addr")
        .context("missing addr field")?
        .as_str()
        .context("fail to get addr")?
        .to_string();

    let clicks = match data.get("clicks") {
        Some(val) => val.as_i64(),
        _ => None,
    }
    .unwrap_or(1i64);

    let tx = match tx_hash {
        Some(_) => {
            sdk::log_info("send tx");
            1
        }
        _ => {
            sdk::log_error("fail to sent tx");
            0
        }
    };

    let mut query: types::DBQuery;
    if is_addr_exist(&addr) {
        query = types::DBQuery {
            statement: String::from(
                "UPDATE Click2Earn SET TotalClicks = TotalClicks + ?, TotalTx = TotalTx + ? WHERE Addr = ?;",
            ),
            params: vec![
                types::Param {
                    int64: Some(clicks),
                    ..Default::default()
                },
                types::Param {
                    int32: Some(tx),
                    ..Default::default()
                },
                types::Param {
                    string: Some(addr),
                    ..Default::default()
                },
            ],
        };
    } else {
        query = types::DBQuery {
            statement: String::from(
                "INSERT INTO Click2Earn (Addr, TotalClicks, TotalTx) VALUES (?, ?, ?);",
            ),
            params: vec![
                types::Param {
                    string: Some(addr),
                    ..Default::default()
                },
                types::Param {
                    int64: Some(clicks),
                    ..Default::default()
                },
                types::Param {
                    int32: Some(tx),
                    ..Default::default()
                },
            ],
        };
    }

    let str = serde_json::to_string(&query).unwrap();
    sdk::set_db(&str)
}

fn is_addr_exist(addr: &String) -> bool {
    let query = types::DBQuery {
        statement: String::from("SELECT 1 FROM Click2Earn WHERE Addr =?;"),
        params: vec![types::Param {
            string: Some(addr.clone()),
            ..Default::default()
        }],
    };
    let str = serde_json::to_string(&query).unwrap();
    let val = match sdk::get_db(&str) {
        Some(val) => val,
        _ => {
            return false;
        }
    };
    val.len() > 0
}
