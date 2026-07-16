ALTER TABLE rides
DROP COLUMN ride_otp_hash,

DROP COLUMN ride_otp_generated_at,

DROP COLUMN ride_otp_verified,

DROP COLUMN ride_otp_attempts;