CREATE TABLE tagged (
  joiner_key TEXT   NOT NULL,
  id         TEXT   NOT NULL,
  tag        TEXT   NOT NULL,
  relation   TEXT   NOT NULL,
  owner_key  TEXT   NOT NULL,
  viewer_key TEXT   NOT NULL

);

CREATE UNIQUE INDEX idx_tagged_uniq ON tagged(owner_key, joiner_key, id, tag);

CREATE        INDEX idx_tagged_tag  ON tagged(viewer_key, tag);

-------------------------

CREATE TABLE tags (
  tag         TEXT    NOT NULL,
  is_internal INTEGER NOT NULL,
  parted_size INTEGER NOT NULL
);

CREATE UNIQUE INDEX idx_tags_uniq ON tags(tag);

CREATE        INDEX idx_tags_int  ON tags(is_internal);
