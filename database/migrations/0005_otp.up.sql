ALTER TABLE rides
DROP COLUMN started_otp;

ALTER TABLE rides
ADD COLUMN started_otp TEXT;