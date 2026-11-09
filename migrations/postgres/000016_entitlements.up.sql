CREATE TABLE IF NOT EXISTS entitlements(
    entitlement_id UUID DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    item_id UUID NOT NULL,
    source TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    obtained_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_till TIMESTAMP,

    PRIMARY KEY(user_id, item_id),

    FOREIGN KEY(user_id) REFERENCES users(id)
    ON DELETE CASCADE,

    FOREIGN KEY(item_id) REFERENCES player_items(id)
    ON DELETE CASCADE
);

CREATE INDEX idx_entitlements_id ON entitlements(entitlement_id);

-- INSERT seed entitlements
INSERT INTO entitlements(entitlement_id, user_id, item_id, source, is_active, obtained_on, valid_till)
VALUES
('2d0b5833-df65-4503-bd4e-721fb89e648c', '5cc59547-508d-4e74-bb36-675eac84e10b', 'b5a76e53-cfa1-4a9f-9749-4436280f5519', 'purchase', true, CURRENT_TIMESTAMP, NULL),
('1fab83ff-5d53-4e5c-8e07-bf608902fc21', '5cc59547-508d-4e74-bb36-675eac84e10b', '1d9fe0fb-c42c-4686-b224-52653727150f', 'purchase', true, CURRENT_TIMESTAMP, NULL),
('f1bbfe3b-e367-4c9a-b995-14d63235c0da', 'dcfd8b0c-fb72-48ed-9cde-ca4285d3f598', 'b5a76e53-cfa1-4a9f-9749-4436280f5519', 'purchase', true, CURRENT_TIMESTAMP, NULL),
('4ccba445-a0a7-41f1-bc6b-228937fb6f4b', 'dcfd8b0c-fb72-48ed-9cde-ca4285d3f598', '1d9fe0fb-c42c-4686-b224-52653727150f', 'purchase', true, CURRENT_TIMESTAMP, NULL);
