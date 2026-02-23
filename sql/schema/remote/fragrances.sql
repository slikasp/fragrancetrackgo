CREATE TABLE public.fragrances (
  id bigserial PRIMARY KEY,

  url          text,
  name         text,
  brand        text,
  country      text,
  gender       text,

  rating_value numeric(3,2),
  rating_count integer,
  year         integer,

  top_notes    text,
  middle_notes text,
  base_notes   text,

  perfumer1    text,
  perfumer2    text,

  accord1      text,
  accord2      text,
  accord3      text,
  accord4      text,
  accord5      text,

  fragrantica_id integer
);