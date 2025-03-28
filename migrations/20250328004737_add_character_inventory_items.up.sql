CREATE TABLE character_inventory_items (
  character_id UUID NOT NULL,
  item_id INTEGER NOT NULL,
  slot INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  PRIMARY KEY (character_id, item_id, slot)
);

