CREATE TABLE provinces (
    code       VARCHAR(10)  PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_provinces_name ON provinces (name);
CREATE INDEX idx_provinces_is_active ON provinces (is_active);

CREATE TABLE cities (
    code          VARCHAR(10)  PRIMARY KEY,
    province_code VARCHAR(10)  NOT NULL REFERENCES provinces (code) ON UPDATE CASCADE ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_cities_province_code ON cities (province_code);
CREATE INDEX idx_cities_name ON cities (name);
CREATE INDEX idx_cities_is_active ON cities (is_active);

CREATE TABLE districts (
    code       VARCHAR(15)  PRIMARY KEY,
    city_code  VARCHAR(10)  NOT NULL REFERENCES cities (code) ON UPDATE CASCADE ON DELETE CASCADE,
    name       VARCHAR(100) NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_districts_city_code ON districts (city_code);
CREATE INDEX idx_districts_name ON districts (name);
CREATE INDEX idx_districts_is_active ON districts (is_active);

CREATE TABLE villages (
    code          VARCHAR(20)  PRIMARY KEY,
    district_code VARCHAR(15)  NOT NULL REFERENCES districts (code) ON UPDATE CASCADE ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    postal_code   VARCHAR(10)  NOT NULL DEFAULT '',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_villages_district_code ON villages (district_code);
CREATE INDEX idx_villages_postal_code ON villages (postal_code);
CREATE INDEX idx_villages_name ON villages (name);
CREATE INDEX idx_villages_is_active ON villages (is_active);
