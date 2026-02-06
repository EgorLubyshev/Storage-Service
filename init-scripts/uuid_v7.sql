CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS uuid AS $$
DECLARE
    unix_ms BIGINT;
    rand_bytes BYTEA;
    uuid_bytes BYTEA;
BEGIN
    -- Получаем текущее время в миллисекундах с эпохи Unix
    unix_ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;

    -- Генерируем 6 случайных байт (48 бит)
    rand_bytes := gen_random_bytes(6);

    -- Формируем 16-байтовый UUIDv7:
    -- 6 байт: timestamp (млс)
    -- 6 байт: random
    -- 4 байта: random (но с установленными битами версии и варианта)
    uuid_bytes := (
        -- Первые 6 байт: timestamp (big-endian)
        (unix_ms >> 16)::BYTEA ||
        (unix_ms & 65535)::BYTEA
    ) || rand_bytes;

    -- Устанавливаем версию (7) в байте 7 (индекс 6)
    -- UUIDv7: байт 6 (0-индексированный) должен быть 0x70 | (случайный бит 4-0)
    -- Но по RFC: версия — 4 бита в байте 6 (позиция 6), старшие 4 бита = 0111 = 0x7
    uuid_bytes := overlay(uuid_bytes PLACING (set_byte(substring(uuid_bytes from 7 for 1), 0, 0x70)) FROM 7 FOR 1);

    -- Устанавливаем вариант (10xx) — байт 8 (индекс 8) должен быть 0x80 | (случайный бит 6-0)
    -- RFC 9562: вариант = 10xx, т.е. старшие 2 бита = 10, значит байт = 0x80 - 0xBF
    -- Мы просто устанавливаем старший бит в 1, остальные — случайные (уже сгенерированы)
    uuid_bytes := overlay(uuid_bytes PLACING (set_byte(substring(uuid_bytes from 9 for 1), 0, (b'10000000' | (get_byte(rand_bytes, 0) & b'00111111'))) ) FROM 9 FOR 1);

    -- Преобразуем в UUID
    RETURN encode(uuid_bytes, 'hex')::uuid;
END;
$$ LANGUAGE plpgsql VOLATILE;

-- CREATE EXTENSION IF NOT EXISTS pg_uuidv7;
