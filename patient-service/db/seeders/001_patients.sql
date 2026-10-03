INSERT INTO patients (nik, name, birth_date, gender, phone, address) VALUES
    ('3301010101900001', 'Budi Santoso', '1990-01-01', 'male', '081234567890', 'Jl. Mawar No. 1'),
    ('3301010202920002', 'Siti Aminah', '1992-02-02', 'female', '081234567891', 'Jl. Melati No. 2'),
    ('3301010303850003', 'Agus Wijaya', '1985-03-03', 'male', '', '')
ON CONFLICT (nik) DO NOTHING;