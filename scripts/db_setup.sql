CREATE database ddoor_db;

\c ddoor_db

DROP table news;

CREATE TABLE news (
  index serial,
  srce text,
  url text PRIMARY KEY,
  description text,
  pub_date DATE NOT NULL DEFAULT CURRENT_DATE,
  ins_date DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE finance_news (
  index serial,
  srce text,
  url text PRIMARY KEY,
  description text,
  pub_date DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ins_date DATE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE news_srce (
  srce text,
  cat text,
  sub_cat text,
  url text PRIMARY KEY,
  tab_name text NOT null,
  ins_date DATE NOT NULL DEFAULT CURRENT_TIMESTAMP
  );


insert into news_srce values
('seekingalpha.com','Finance','General','https://seekingalpha.com/sitemap_news.xml','finance_news');

insert into news_srce values
('ycombinator','News','General','https://news.ycombinator.com/rss','finance_news');

insert into news_srce values
('wpost','News','Finance','https://www.washingtonpost.com/news-sitemaps/business.xml','finance_news');

insert into news_srce values
('bbc','world','General','http://feeds.bbci.co.uk/news/world/rss.xml','finance_news');

insert into news_srce values
('yahoo','world','General','https://www.yahoo.com/news/rss/world','finance_news');
