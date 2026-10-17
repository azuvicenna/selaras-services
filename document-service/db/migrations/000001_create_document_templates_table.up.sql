CREATE TABLE document_templates (
    id            CHAR(26)     PRIMARY KEY,
    code          VARCHAR(50)  NOT NULL UNIQUE,
    name          VARCHAR(150) NOT NULL,
    type          SMALLINT     NOT NULL DEFAULT 1,
    typst_content TEXT         NOT NULL DEFAULT '',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    version       INT          NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type ON document_templates (type);
CREATE INDEX idx_document_templates_is_active ON document_templates (is_active);
