BEGIN TRANSACTION;
DROP TABLE IF EXISTS `playsearch`;
CREATE VIRTUAL TABLE playsearch USING fts5(playsrowid, text);
--please remember that CREATE VIRTUAL TABLE playsearch USING fts5
-- creatses all tables needed to store addition info for the search
-- After the sql exmport, manually removed the creation of: playsearch_idx, playsearch_docsize, playsearch_data, playsearch_content, playsearch_config
-- Or do not export those tables.
-- Otherwise the import of new db fil failed with:
--      error creating shadow table playsearch_data: table 'playsearch_data' already exists

DROP TABLE IF EXISTS `Playlist`;
CREATE TABLE IF NOT EXISTS `Playlist` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`Name`	TEXT UNIQUE
);

DROP TABLE IF EXISTS `Item`;
CREATE TABLE IF NOT EXISTS `Item` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`URI`	TEXT,
	`Info`	TEXT,
	`ItemType`	INTEGER,
	`Description`	TEXT,
	`MetaTitle`	TEXT,
	`MetaFileType`	TEXT,
	`MetaAlbum`	TEXT,
	`MetaArtist`	TEXT,
	`MetaAlbumArtist`	TEXT
);
DROP TABLE IF EXISTS `History`;
CREATE TABLE IF NOT EXISTS `History` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`Timestamp`	INTEGER,
	`URI`	TEXT,
	`Title`	TEXT,
	`Description`	TEXT,
	`Duration`	TEXT,
	`PlayPosition`	INTEGER,
	`DurationInSec`	INTEGER,
	`Type`	TEXT
);
DROP TABLE IF EXISTS `Current`;
CREATE TABLE IF NOT EXISTS `Current` (
	`id`	INTEGER,
	`ListName`	TEXT NOT NULL,
	`Volatile`	INTEGER,
	`URI`	TEXT,
	`Info`	TEXT,
	`ItemType`	INTEGER,
	PRIMARY KEY(`id`)
);
COMMIT;
