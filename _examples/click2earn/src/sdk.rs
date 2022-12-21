use anyhow::{bail, Result};

#[link(wasm_import_module = "env")]
extern "C" {
    fn ws_log(log_level: i32, ptr: *const u8, size: i32) -> i32;
    fn ws_get_sql_db(
        key_ptr: *const u8,
        key_size: i32,
        return_ptr: *const *mut u8,
        return_size: *const i32,
    ) -> i32;
    fn ws_set_sql_db(ptr: *const u8, size: i32) -> i32;
    fn ws_get_data(resource_id: i32, return_ptr: *const *mut u8, return_size: *const i32) -> i32;

    fn ws_send_tx(
        payload_ptr: *const u8,
        payload_size: i32,
        return_hash_ptr: *const *mut u8,
        return_hash_size: *const i32,
    ) -> i32;
}

enum LogLevel {
    Trace = 1,
    Debug,
    Info,
    Warn,
    Error,
}

pub fn log_info(str: &str) {
    unsafe { ws_log(LogLevel::Info as _, str.as_bytes().as_ptr(), str.len() as _) };
}

pub fn log_error(str: &str) {
    unsafe {
        ws_log(
            LogLevel::Error as _,
            str.as_bytes().as_ptr(),
            str.len() as _,
        )
    };
}

pub fn set_db(encoded: &String) -> Result<()> {
    match unsafe { ws_set_sql_db(encoded.as_ptr(), encoded.len() as _) } {
        0 => Ok(()),
        _ => bail!("fail to set db"),
    }
}

pub fn get_db(encoded: &String) -> Option<Vec<u8>> {
    let data_ptr = &mut (0 as i32) as *const _ as *const *mut u8;
    let data_size = &(0 as i32);
    match unsafe { ws_get_sql_db(encoded.as_ptr(), encoded.len() as _, data_ptr, data_size) } {
        0 => Some(unsafe { Vec::from_raw_parts(*data_ptr, *data_size as _, *data_size as _) }),
        _ => None,
    }
}

pub fn get_data(resource_id: i32) -> Option<Vec<u8>> {
    let data_ptr = &mut (0 as i32) as *const _ as *const *mut u8;
    let data_size = &(0 as i32);
    match unsafe { ws_get_data(resource_id, data_ptr, data_size) } {
        0 => Some(unsafe { Vec::from_raw_parts(*data_ptr, *data_size as _, *data_size as _) }),
        _ => None,
    }
}

pub fn send_tx(to: &String, value: &String, data: &String) -> Result<String> {
    let tx = crate::types::Tx {
        to: to.clone(),
        value: value.clone(),
        data: data.clone(),
    };
    let str = serde_json::to_string(&tx)?;
    let data_ptr = &mut (0 as i32) as *const _ as *const *mut u8;
    let data_size = &(0 as i32);
    match unsafe { ws_send_tx(str.as_ptr(), str.len() as _, data_ptr, data_size) } {
        0 => Ok(unsafe { String::from_raw_parts(*data_ptr, *data_size as _, *data_size as _) }),
        _ => bail!("fail to send tx"),
    }
}
