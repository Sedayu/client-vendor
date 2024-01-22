-- Create Vendors Table
CREATE TABLE vendors (
     id SERIAL PRIMARY KEY,
     name VARCHAR(255) NOT NULL,
     email VARCHAR(255) NOT NULL,
     phone_number VARCHAR(20),
     address TEXT,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL

);

-- Create an index on the email column for faster lookups (optional)
CREATE INDEX idx_vendor_email ON vendors(email);