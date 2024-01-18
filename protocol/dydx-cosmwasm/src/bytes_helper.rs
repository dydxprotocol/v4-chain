use num_bigint::{BigInt, Sign};

pub(crate) fn bytes_to_bigint(bytes: Vec<u8>) -> BigInt {
    if bytes.len() <= 1 {
        return BigInt::from(0);
    }

    let sign = if bytes[0] & 1 == 1 { Sign::Minus } else { Sign::Plus };
    let abs = BigInt::from_bytes_be(sign, &bytes[1..]);

    abs
}

#[cfg(test)]
mod tests {
    use super::*;
    use num_bigint::BigInt;

    #[test]
    fn test_bytes_to_bigint() {
        let test_cases = vec![
            (BigInt::from(0), vec![0x02]),
            (BigInt::from(0), vec![0x02]), // -0 is same as 0
            (BigInt::from(1), vec![0x02, 0x01]),
            (BigInt::from(-1), vec![0x03, 0x01]),
            (BigInt::from(255), vec![0x02, 0xFF]),
            (BigInt::from(-255), vec![0x03, 0xFF]),
            (BigInt::from(256), vec![0x02, 0x01, 0x00]),
            (BigInt::from(-256), vec![0x03, 0x01, 0x00]),
            (BigInt::from(123456789), vec![0x02, 0x07, 0x5b, 0xcd, 0x15]),
            (BigInt::from(-123456789), vec![0x03, 0x07, 0x5b, 0xcd, 0x15]),
            (BigInt::from(123456789123456789i64), vec![0x02, 0x01, 0xb6, 0x9b, 0x4b, 0xac, 0xd0, 0x5f, 0x15]),
            (BigInt::from(-123456789123456789i64), vec![0x03, 0x01, 0xb6, 0x9b, 0x4b, 0xac, 0xd0, 0x5f, 0x15]),
            (BigInt::parse_bytes(b"123456789123456789123456789", 10).unwrap(), vec![0x02, 0x66, 0x1e, 0xfd, 0xf2, 0xe3, 0xb1, 0x9f, 0x7c, 0x04, 0x5f, 0x15]),
            (BigInt::parse_bytes(b"-123456789123456789123456789", 10).unwrap(), vec![0x03, 0x66, 0x1e, 0xfd, 0xf2, 0xe3, 0xb1, 0x9f, 0x7c, 0x04, 0x5f, 0x15]),
        ];

        for (expected, bytes) in test_cases {
            assert_eq!(bytes_to_bigint(bytes), expected);
        }
    }
}
