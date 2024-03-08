CREATE TABLE IF NOT EXISTS public."users"
(
    id bigserial NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(20) NOT NULL DEFAULT 'USER'::character varying,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    created_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    modified_at timestamp without time zone NOT NULL DEFAULT now(),
    modified_by character varying(255) NOT NULL DEFAULT 'SYSTEM'::character varying,
    deleted_at timestamp without time zone DEFAULT NULL,
    deleted_by character varying(255) DEFAULT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
);