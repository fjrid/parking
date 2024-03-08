CREATE TABLE IF NOT EXISTS public."parking"
(
    id bigserial NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp without time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp without time zone DEFAULT NULL,
    deleted_by character varying(255) DEFAULT NULL,
    CONSTRAINT parking_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."parking_slot"
(
    id bigserial NOT NULL,
	parking_id bigint NOT NULL,
    number bigint NOT NULL,
    status character varying(20) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp without time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp without time zone DEFAULT NULL,
    deleted_by character varying(255) DEFAULT NULL,
    CONSTRAINT parking_lot_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."parking_book"
(
    id bigserial NOT NULL,
    user_id bigint NOT NULL,
	parking_slot_id bigint NOT NULL,
    start_time timestamp without time zone NOT NULL DEFAULT now(),
    end_time timestamp without time zone DEFAULT NULL,
    fee numeric DEFAULT 0,
    status character varying(20) NOT NULL,
	car_number character varying(20) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp without time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp without time zone DEFAULT NULL,
    deleted_by character varying(255) DEFAULT NULL,
    CONSTRAINT parking_book_pkey PRIMARY KEY (id)
);