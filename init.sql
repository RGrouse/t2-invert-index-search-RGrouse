CREATE TABLE public.words
(
  id serial NOT NULL,
  word text,
  CONSTRAINT words_pk PRIMARY KEY (id)
);


CREATE TABLE public.sources
(
  id serial NOT NULL,
  source text,
  CONSTRAINT sources_pk PRIMARY KEY (id)
);


CREATE TABLE public.index
(
  id serial NOT NULL,
  word_id integer,
  source_id integer,
  weight integer,
  CONSTRAINT index_pk PRIMARY KEY (id),
  CONSTRAINT index_word_id FOREIGN KEY (word_id)
      REFERENCES public.words (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT index_source_id FOREIGN KEY (source_id)
      REFERENCES public.sources (id) MATCH SIMPLE
      ON UPDATE NO ACTION ON DELETE CASCADE
);