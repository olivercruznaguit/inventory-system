CREATE TABLE stock_movements (
    id                  BIGSERIAL PRIMARY KEY,
    product_id          BIGINT NOT NULL,
    type                VARCHAR(10) NOT NULL,
    quantity            INT NOT NULL,
    remaining_quantity  INT NOT NULL,
    reason              TEXT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT fk_stock_movements_product FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,

    CONSTRAINT chk_stock_movements_type CHECK (type IN ('IN', 'OUT')),

    CONSTRAINT chk_stock_movements_quantity CHECK (quantity > 0),

    CONSTRAINT chk_stock_movements_remaining CHECK (remaining_quantity >= 0)
);

-- Index
CREATE INDEX idx_stock_movements_product_created_at ON stock_movements(product_id, created_at DESC);