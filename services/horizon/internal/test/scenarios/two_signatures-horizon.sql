--
-- PostgreSQL database dump
--

-- Dumped from database version 11.2 (Ubuntu 11.2-1.pgdg18.04+1)
-- Dumped by pg_dump version 11.2 (Ubuntu 11.2-1.pgdg18.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;

ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_counter_asset_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_counter_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_base_asset_id_fkey;
ALTER TABLE IF EXISTS ONLY public.history_trades DROP CONSTRAINT IF EXISTS history_trades_base_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.asset_stats DROP CONSTRAINT IF EXISTS asset_stats_id_fkey;
DROP INDEX IF EXISTS public.trade_effects_by_order_book;
DROP INDEX IF EXISTS public.index_history_transactions_on_id;
DROP INDEX IF EXISTS public.index_history_operations_on_type;
DROP INDEX IF EXISTS public.index_history_operations_on_transaction_id;
DROP INDEX IF EXISTS public.index_history_operations_on_id;
DROP INDEX IF EXISTS public.index_history_ledgers_on_sequence;
DROP INDEX IF EXISTS public.index_history_ledgers_on_previous_ledger_hash;
DROP INDEX IF EXISTS public.index_history_ledgers_on_ledger_hash;
DROP INDEX IF EXISTS public.index_history_ledgers_on_importer_version;
DROP INDEX IF EXISTS public.index_history_ledgers_on_id;
DROP INDEX IF EXISTS public.index_history_ledgers_on_closed_at;
DROP INDEX IF EXISTS public.index_history_effects_on_type;
DROP INDEX IF EXISTS public.index_history_accounts_on_id;
DROP INDEX IF EXISTS public.index_history_accounts_on_address;
DROP INDEX IF EXISTS public.htrd_time_lookup;
DROP INDEX IF EXISTS public.htrd_pid;
DROP INDEX IF EXISTS public.htrd_pair_time_lookup;
DROP INDEX IF EXISTS public.htrd_counter_lookup;
DROP INDEX IF EXISTS public.htrd_by_offer;
DROP INDEX IF EXISTS public.htrd_by_counter_offer;
DROP INDEX IF EXISTS public.htrd_by_counter_account;
DROP INDEX IF EXISTS public.htrd_by_base_offer;
DROP INDEX IF EXISTS public.htrd_by_base_account;
DROP INDEX IF EXISTS public.htp_by_htid;
DROP INDEX IF EXISTS public.hs_transaction_by_id;
DROP INDEX IF EXISTS public.hs_ledger_by_id;
DROP INDEX IF EXISTS public.hop_by_hoid;
DROP INDEX IF EXISTS public.hist_tx_p_id;
DROP INDEX IF EXISTS public.hist_op_p_id;
DROP INDEX IF EXISTS public.hist_e_id;
DROP INDEX IF EXISTS public.hist_e_by_order;
DROP INDEX IF EXISTS public.by_ledger;
DROP INDEX IF EXISTS public.by_hash;
DROP INDEX IF EXISTS public.by_account;
DROP INDEX IF EXISTS public.asset_by_issuer;
DROP INDEX IF EXISTS public.asset_by_code;
ALTER TABLE IF EXISTS ONLY public.history_transaction_participants DROP CONSTRAINT IF EXISTS history_transaction_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_operation_participants DROP CONSTRAINT IF EXISTS history_operation_participants_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_pkey;
ALTER TABLE IF EXISTS ONLY public.history_assets DROP CONSTRAINT IF EXISTS history_assets_asset_code_asset_type_asset_issuer_key;
ALTER TABLE IF EXISTS ONLY public.gorp_migrations DROP CONSTRAINT IF EXISTS gorp_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.asset_stats DROP CONSTRAINT IF EXISTS asset_stats_pkey;
ALTER TABLE IF EXISTS public.history_transaction_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_operation_participants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.history_assets ALTER COLUMN id DROP DEFAULT;
DROP TABLE IF EXISTS public.history_transactions;
DROP SEQUENCE IF EXISTS public.history_transaction_participants_id_seq;
DROP TABLE IF EXISTS public.history_transaction_participants;
DROP TABLE IF EXISTS public.history_trades;
DROP TABLE IF EXISTS public.history_operations;
DROP SEQUENCE IF EXISTS public.history_operation_participants_id_seq;
DROP TABLE IF EXISTS public.history_operation_participants;
DROP TABLE IF EXISTS public.history_ledgers;
DROP TABLE IF EXISTS public.history_effects;
DROP SEQUENCE IF EXISTS public.history_assets_id_seq;
DROP TABLE IF EXISTS public.history_assets;
DROP TABLE IF EXISTS public.history_accounts;
DROP SEQUENCE IF EXISTS public.history_accounts_id_seq;
DROP TABLE IF EXISTS public.gorp_migrations;
DROP TABLE IF EXISTS public.asset_stats;
DROP AGGREGATE IF EXISTS public.min_price(numeric[]);
DROP AGGREGATE IF EXISTS public.max_price(numeric[]);
DROP AGGREGATE IF EXISTS public.last(anyelement);
DROP AGGREGATE IF EXISTS public.first(anyelement);
DROP FUNCTION IF EXISTS public.min_price_agg(numeric[], numeric[]);
DROP FUNCTION IF EXISTS public.max_price_agg(numeric[], numeric[]);
DROP FUNCTION IF EXISTS public.last_agg(anyelement, anyelement);
DROP FUNCTION IF EXISTS public.first_agg(anyelement, anyelement);
--
-- Name: first_agg(anyelement, anyelement); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.first_agg(anyelement, anyelement) RETURNS anyelement
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT $1 $_$;


--
-- Name: last_agg(anyelement, anyelement); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.last_agg(anyelement, anyelement) RETURNS anyelement
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT $2 $_$;


--
-- Name: max_price_agg(numeric[], numeric[]); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.max_price_agg(numeric[], numeric[]) RETURNS numeric[]
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT (
  CASE WHEN $1[1]/$1[2]>$2[1]/$2[2] THEN $1 ELSE $2 END) $_$;


--
-- Name: min_price_agg(numeric[], numeric[]); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.min_price_agg(numeric[], numeric[]) RETURNS numeric[]
    LANGUAGE sql IMMUTABLE STRICT
    AS $_$ SELECT (
  CASE WHEN $1[1]/$1[2]<$2[1]/$2[2] THEN $1 ELSE $2 END) $_$;


--
-- Name: first(anyelement); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE public.first(anyelement) (
    SFUNC = public.first_agg,
    STYPE = anyelement
);


--
-- Name: last(anyelement); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE public.last(anyelement) (
    SFUNC = public.last_agg,
    STYPE = anyelement
);


--
-- Name: max_price(numeric[]); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE public.max_price(numeric[]) (
    SFUNC = public.max_price_agg,
    STYPE = numeric[]
);


--
-- Name: min_price(numeric[]); Type: AGGREGATE; Schema: public; Owner: -
--

CREATE AGGREGATE public.min_price(numeric[]) (
    SFUNC = public.min_price_agg,
    STYPE = numeric[]
);


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: asset_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.asset_stats (
    id bigint NOT NULL,
    amount character varying NOT NULL,
    num_accounts integer NOT NULL,
    flags smallint NOT NULL,
    toml character varying(255) NOT NULL
);


--
-- Name: gorp_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.gorp_migrations (
    id text NOT NULL,
    applied_at timestamp with time zone
);


--
-- Name: history_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.history_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_accounts (
    id bigint DEFAULT nextval('public.history_accounts_id_seq'::regclass) NOT NULL,
    address character varying(64)
);


--
-- Name: history_assets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_assets (
    id integer NOT NULL,
    asset_type character varying(64) NOT NULL,
    asset_code character varying(12) NOT NULL,
    asset_issuer character varying(56) NOT NULL
);


--
-- Name: history_assets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.history_assets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_assets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.history_assets_id_seq OWNED BY public.history_assets.id;


--
-- Name: history_effects; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_effects (
    history_account_id bigint NOT NULL,
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    type integer NOT NULL,
    details jsonb
);


--
-- Name: history_ledgers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_ledgers (
    sequence integer NOT NULL,
    ledger_hash character varying(64) NOT NULL,
    previous_ledger_hash character varying(64),
    transaction_count integer DEFAULT 0 NOT NULL,
    operation_count integer DEFAULT 0 NOT NULL,
    closed_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    id bigint,
    importer_version integer DEFAULT 1 NOT NULL,
    total_coins bigint NOT NULL,
    fee_pool bigint NOT NULL,
    base_fee integer NOT NULL,
    base_reserve integer NOT NULL,
    max_tx_set_size integer NOT NULL,
    protocol_version integer DEFAULT 0 NOT NULL,
    ledger_header text,
    successful_transaction_count integer,
    failed_transaction_count integer
);


--
-- Name: history_operation_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_operation_participants (
    id integer NOT NULL,
    history_operation_id bigint NOT NULL,
    history_account_id bigint NOT NULL
);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.history_operation_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.history_operation_participants_id_seq OWNED BY public.history_operation_participants.id;


--
-- Name: history_operations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_operations (
    id bigint NOT NULL,
    transaction_id bigint NOT NULL,
    application_order integer NOT NULL,
    type integer NOT NULL,
    details jsonb,
    source_account character varying(64) DEFAULT ''::character varying NOT NULL
);


--
-- Name: history_trades; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_trades (
    history_operation_id bigint NOT NULL,
    "order" integer NOT NULL,
    ledger_closed_at timestamp without time zone NOT NULL,
    offer_id bigint NOT NULL,
    base_account_id bigint NOT NULL,
    base_asset_id bigint NOT NULL,
    base_amount bigint NOT NULL,
    counter_account_id bigint NOT NULL,
    counter_asset_id bigint NOT NULL,
    counter_amount bigint NOT NULL,
    base_is_seller boolean,
    price_n bigint,
    price_d bigint,
    base_offer_id bigint,
    counter_offer_id bigint,
    CONSTRAINT history_trades_base_amount_check CHECK ((base_amount > 0)),
    CONSTRAINT history_trades_check CHECK ((base_asset_id < counter_asset_id)),
    CONSTRAINT history_trades_counter_amount_check CHECK ((counter_amount > 0))
);


--
-- Name: history_transaction_participants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_transaction_participants (
    id integer NOT NULL,
    history_transaction_id bigint NOT NULL,
    history_account_id bigint NOT NULL
);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.history_transaction_participants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.history_transaction_participants_id_seq OWNED BY public.history_transaction_participants.id;


--
-- Name: history_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.history_transactions (
    transaction_hash character varying(64) NOT NULL,
    ledger_sequence integer NOT NULL,
    application_order integer NOT NULL,
    account character varying(64) NOT NULL,
    account_sequence bigint NOT NULL,
    fee_paid integer NOT NULL,
    operation_count integer NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    id bigint,
    tx_envelope text NOT NULL,
    tx_result text NOT NULL,
    tx_meta text NOT NULL,
    tx_fee_meta text NOT NULL,
    signatures character varying(96)[] DEFAULT '{}'::character varying[] NOT NULL,
    memo_type character varying DEFAULT 'none'::character varying NOT NULL,
    memo character varying,
    time_bounds int8range
);


--
-- Name: history_assets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_assets ALTER COLUMN id SET DEFAULT nextval('public.history_assets_id_seq'::regclass);


--
-- Name: history_operation_participants id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_operation_participants ALTER COLUMN id SET DEFAULT nextval('public.history_operation_participants_id_seq'::regclass);


--
-- Name: history_transaction_participants id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_transaction_participants ALTER COLUMN id SET DEFAULT nextval('public.history_transaction_participants_id_seq'::regclass);


--
-- Data for Name: asset_stats; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: gorp_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.gorp_migrations VALUES ('1_initial_schema.sql', '2019-01-31 19:27:26.713477+02');
INSERT INTO public.gorp_migrations VALUES ('2_index_participants_by_toid.sql', '2019-01-31 19:27:26.726446+02');
INSERT INTO public.gorp_migrations VALUES ('3_use_sequence_in_history_accounts.sql', '2019-01-31 19:27:26.730946+02');
INSERT INTO public.gorp_migrations VALUES ('4_add_protocol_version.sql', '2019-01-31 19:27:26.741471+02');
INSERT INTO public.gorp_migrations VALUES ('5_create_trades_table.sql', '2019-01-31 19:27:26.755266+02');
INSERT INTO public.gorp_migrations VALUES ('6_create_assets_table.sql', '2019-01-31 19:27:26.761602+02');
INSERT INTO public.gorp_migrations VALUES ('7_modify_trades_table.sql', '2019-01-31 19:27:26.775197+02');
INSERT INTO public.gorp_migrations VALUES ('8_create_asset_stats_table.sql', '2019-01-31 19:27:26.782053+02');
INSERT INTO public.gorp_migrations VALUES ('8_add_aggregators.sql', '2019-01-31 19:27:26.784662+02');
INSERT INTO public.gorp_migrations VALUES ('9_add_header_xdr.sql', '2019-01-31 19:27:26.787483+02');
INSERT INTO public.gorp_migrations VALUES ('10_add_trades_price.sql', '2019-01-31 19:27:26.790524+02');
INSERT INTO public.gorp_migrations VALUES ('11_add_trades_account_index.sql', '2019-01-31 19:27:26.795202+02');
INSERT INTO public.gorp_migrations VALUES ('12_asset_stats_amount_string.sql', '2019-01-31 19:27:26.801587+02');
INSERT INTO public.gorp_migrations VALUES ('13_trade_offer_ids.sql', '2019-01-31 19:27:26.809634+02');
INSERT INTO public.gorp_migrations VALUES ('14_fix_asset_toml_field.sql', '2019-01-31 19:27:26.811268+02');
INSERT INTO public.gorp_migrations VALUES ('15_ledger_failed_txs.sql', '2019-01-31 19:27:26.81283+02');


--
-- Data for Name: history_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_accounts VALUES (1, 'GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT');
INSERT INTO public.history_accounts VALUES (2, 'GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN');
INSERT INTO public.history_accounts VALUES (3, 'GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2');
INSERT INTO public.history_accounts VALUES (4, 'GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL');
INSERT INTO public.history_accounts VALUES (5, 'GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ');
INSERT INTO public.history_accounts VALUES (6, 'GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7');
INSERT INTO public.history_accounts VALUES (7, 'GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ');
INSERT INTO public.history_accounts VALUES (8, 'GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON');
INSERT INTO public.history_accounts VALUES (9, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO public.history_accounts VALUES (10, 'GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y');
INSERT INTO public.history_accounts VALUES (11, 'GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6');
INSERT INTO public.history_accounts VALUES (12, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_accounts VALUES (13, 'GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2');


--
-- Data for Name: history_assets; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: history_effects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_effects VALUES (1, 21474840577, 1, 12, '{"weight": 1, "public_key": "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT"}');
INSERT INTO public.history_effects VALUES (1, 21474840577, 2, 10, '{"weight": 1, "public_key": "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN"}');
INSERT INTO public.history_effects VALUES (2, 21474844673, 1, 12, '{"weight": 1, "public_key": "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN"}');
INSERT INTO public.history_effects VALUES (2, 21474844673, 2, 10, '{"weight": 1, "public_key": "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT"}');
INSERT INTO public.history_effects VALUES (3, 17179873281, 1, 12, '{"weight": 1, "public_key": "GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2"}');
INSERT INTO public.history_effects VALUES (3, 17179873281, 2, 10, '{"weight": 1, "public_key": "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ"}');
INSERT INTO public.history_effects VALUES (4, 17179877377, 1, 12, '{"weight": 1, "public_key": "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL"}');
INSERT INTO public.history_effects VALUES (4, 17179877377, 2, 10, '{"weight": 1, "public_key": "GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2"}');
INSERT INTO public.history_effects VALUES (5, 17179881473, 1, 12, '{"weight": 1, "public_key": "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ"}');
INSERT INTO public.history_effects VALUES (5, 17179881473, 2, 10, '{"weight": 1, "public_key": "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL"}');
INSERT INTO public.history_effects VALUES (6, 12884905985, 1, 12, '{"weight": 1, "public_key": "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7"}');
INSERT INTO public.history_effects VALUES (6, 12884905985, 2, 10, '{"weight": 1, "public_key": "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y"}');
INSERT INTO public.history_effects VALUES (7, 12884910081, 1, 12, '{"weight": 1, "public_key": "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ"}');
INSERT INTO public.history_effects VALUES (7, 12884910081, 2, 10, '{"weight": 1, "public_key": "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6"}');
INSERT INTO public.history_effects VALUES (8, 12884914177, 1, 2, '{"amount": "500.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (9, 12884914177, 2, 3, '{"amount": "500.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (10, 12884918273, 1, 12, '{"weight": 1, "public_key": "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y"}');
INSERT INTO public.history_effects VALUES (10, 12884918273, 2, 10, '{"weight": 1, "public_key": "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ"}');
INSERT INTO public.history_effects VALUES (9, 12884922369, 1, 12, '{"weight": 1, "public_key": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO public.history_effects VALUES (9, 12884922369, 2, 10, '{"weight": 1, "public_key": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON"}');
INSERT INTO public.history_effects VALUES (9, 8589938689, 1, 0, '{"starting_balance": "10000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589938689, 2, 3, '{"amount": "10000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (9, 8589938689, 3, 10, '{"weight": 1, "public_key": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"}');
INSERT INTO public.history_effects VALUES (13, 8589942785, 1, 0, '{"starting_balance": "10000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589942785, 2, 3, '{"amount": "10000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (13, 8589942785, 3, 10, '{"weight": 1, "public_key": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2"}');
INSERT INTO public.history_effects VALUES (8, 8589946881, 1, 0, '{"starting_balance": "10000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589946881, 2, 3, '{"amount": "10000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (8, 8589946881, 3, 10, '{"weight": 1, "public_key": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON"}');
INSERT INTO public.history_effects VALUES (6, 8589950977, 1, 0, '{"starting_balance": "10000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589950977, 2, 3, '{"amount": "10000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (6, 8589950977, 3, 10, '{"weight": 1, "public_key": "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7"}');
INSERT INTO public.history_effects VALUES (10, 8589955073, 1, 0, '{"starting_balance": "20000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589955073, 2, 3, '{"amount": "20000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (10, 8589955073, 3, 10, '{"weight": 1, "public_key": "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y"}');
INSERT INTO public.history_effects VALUES (7, 8589959169, 1, 0, '{"starting_balance": "30000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589959169, 2, 3, '{"amount": "30000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (7, 8589959169, 3, 10, '{"weight": 1, "public_key": "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ"}');
INSERT INTO public.history_effects VALUES (11, 8589963265, 1, 0, '{"starting_balance": "40000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589963265, 2, 3, '{"amount": "40000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (11, 8589963265, 3, 10, '{"weight": 1, "public_key": "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6"}');
INSERT INTO public.history_effects VALUES (5, 8589967361, 1, 0, '{"starting_balance": "50000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589967361, 2, 3, '{"amount": "50000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (5, 8589967361, 3, 10, '{"weight": 1, "public_key": "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ"}');
INSERT INTO public.history_effects VALUES (3, 8589971457, 1, 0, '{"starting_balance": "60000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589971457, 2, 3, '{"amount": "60000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (3, 8589971457, 3, 10, '{"weight": 1, "public_key": "GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2"}');
INSERT INTO public.history_effects VALUES (4, 8589975553, 1, 0, '{"starting_balance": "70000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589975553, 2, 3, '{"amount": "70000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (4, 8589975553, 3, 10, '{"weight": 1, "public_key": "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL"}');
INSERT INTO public.history_effects VALUES (2, 8589979649, 1, 0, '{"starting_balance": "80000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589979649, 2, 3, '{"amount": "80000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (2, 8589979649, 3, 10, '{"weight": 1, "public_key": "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN"}');
INSERT INTO public.history_effects VALUES (1, 8589983745, 1, 0, '{"starting_balance": "90000.00000"}');
INSERT INTO public.history_effects VALUES (12, 8589983745, 2, 3, '{"amount": "90000.00000", "asset_type": "native"}');
INSERT INTO public.history_effects VALUES (1, 8589983745, 3, 10, '{"weight": 1, "public_key": "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT"}');


--
-- Data for Name: history_ledgers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_ledgers VALUES (5, 'ab91d69af2f5bced4484203e5cd512bc7c5f9270ed75956c328939600e32071c', '64d249d70e6547e3af0ab5bf87a90e3522dae2563c684f34fd720be31218c194', 2, 2, '2019-06-13 14:07:46', '2019-06-13 14:07:46.500982', '2019-06-13 14:07:46.500997', 21474836480, 15, 1000000000000000000, 2200, 100, 100000000, 10000, 9, 'AAAACWTSSdcOZUfjrwq1v4epDjUi2uJWPGhPNP1yC+MSGMGUnCzCrn3GG9Kx1c2O9dhCk0nag8rF8sVXn7s1HL17vyAAAAAAXQJYsgAAAAAAAAAAM4wUQycfXLO+E7jYXg5TRQv5QIpcX1qzGv0MqMjIOdDle1Pk/bi6Poy2M2/PEmtWzGhd5LWep4Gm4/GcY8akGAAAAAUN4Lazp2QAAAAAAAAAAAiYAAAAAAAAAAAAAAAAAAAAZAX14QAAACcQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 2, 0);
INSERT INTO public.history_ledgers VALUES (4, '64d249d70e6547e3af0ab5bf87a90e3522dae2563c684f34fd720be31218c194', '14b54b98a549a56a4d58ceeb5e4ade43091abe8cae71669abee0d2b59b96acc4', 3, 3, '2019-06-13 14:07:45', '2019-06-13 14:07:46.515142', '2019-06-13 14:07:46.515156', 17179869184, 15, 1000000000000000000, 2000, 100, 100000000, 10000, 9, 'AAAACRS1S5ilSaVqTVjO615K3kMJGr6MrnFmmr7g0rWblqzE9txgnvfrUyVUIyPTVMxOj2JYQCyeJ0Egs6ajOiru64kAAAAAXQJYsQAAAAAAAAAAfYEdlm6x8PddC4f0RSpOGXPAHQu40GYwoxti06dAQAY86oBNWB4SjUpEvqANtGyxyV5LiJTBgsvc3D5H6h88HgAAAAQN4Lazp2QAAAAAAAAAAAfQAAAAAAAAAAAAAAAAAAAAZAX14QAAACcQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 3, 0);
INSERT INTO public.history_ledgers VALUES (3, '14b54b98a549a56a4d58ceeb5e4ade43091abe8cae71669abee0d2b59b96acc4', 'bc5df90babfff3abd389fa77875ef5311eaf32e0700812f88cabd1c4bb39a75c', 5, 5, '2019-06-13 14:07:44', '2019-06-13 14:07:46.523613', '2019-06-13 14:07:46.523626', 12884901888, 15, 1000000000000000000, 1700, 100, 100000000, 10000, 9, 'AAAACbxd+Qur//Or04n6d4de9TEerzLgcAgS+Iyr0cS7OadcAKOK0aXWViKV6PryFJuBIAHPP1y8x6iwX6sbBkO6G6AAAAAAXQJYsAAAAAAAAAAAOL/dAce7qKMx3NaIabRXOwXcY3TnsGMlSBqNZbfqtENBln542JjeHDplt1lLjFbvV3ngUm1wWbMnP4EqFg6bewAAAAMN4Lazp2QAAAAAAAAAAAakAAAAAAAAAAAAAAAAAAAAZAX14QAAACcQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 5, 0);
INSERT INTO public.history_ledgers VALUES (2, 'bc5df90babfff3abd389fa77875ef5311eaf32e0700812f88cabd1c4bb39a75c', '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', 12, 12, '2019-06-13 14:07:43', '2019-06-13 14:07:46.536476', '2019-06-13 14:07:46.536489', 8589934592, 15, 1000000000000000000, 1200, 100, 100000000, 10000, 9, 'AAAACWPZj1Nu5o0bJ7W4nyOvUxG3Vpok+vFAOtC1K2M7B76Zn3bJCNfZe3EY48Dr9Xj80NOuIRWWAgwKf0rn4/e1fMgAAAAAXQJYrwAAAAIAAAAIAAAAAQAAAAkAAAAIAAAAAwAAJxAAAAAAhZaqGcW4dpyjg5Try1YR8EtEciWQVGfdyZQy3s5Q1pJH1xCASEyyC4dzuqYaQlTDgoBh10/ZN5xU6eLzJALcwAAAAAIN4Lazp2QAAAAAAAAAAASwAAAAAAAAAAAAAAAAAAAAZAX14QAAACcQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 12, 0);
INSERT INTO public.history_ledgers VALUES (1, '63d98f536ee68d1b27b5b89f23af5311b7569a24faf1403ad0b52b633b07be99', NULL, 0, 0, '1970-01-01 00:00:00', '2019-06-13 14:07:46.553002', '2019-06-13 14:07:46.553022', 4294967296, 15, 1000000000000000000, 0, 100, 100000000, 100, 0, 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABXKi4y/ySKB7DnD9H20xjB+s0gtswIwz1XdSWYaBJaFgAAAAEN4Lazp2QAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAX14QAAAABkAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA', 0, 0);


--
-- Data for Name: history_operation_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_operation_participants VALUES (1, 21474840577, 1);
INSERT INTO public.history_operation_participants VALUES (2, 21474844673, 2);
INSERT INTO public.history_operation_participants VALUES (3, 17179873281, 3);
INSERT INTO public.history_operation_participants VALUES (4, 17179877377, 4);
INSERT INTO public.history_operation_participants VALUES (5, 17179881473, 5);
INSERT INTO public.history_operation_participants VALUES (6, 12884905985, 6);
INSERT INTO public.history_operation_participants VALUES (7, 12884910081, 7);
INSERT INTO public.history_operation_participants VALUES (8, 12884914177, 8);
INSERT INTO public.history_operation_participants VALUES (9, 12884914177, 9);
INSERT INTO public.history_operation_participants VALUES (10, 12884918273, 10);
INSERT INTO public.history_operation_participants VALUES (11, 12884922369, 9);
INSERT INTO public.history_operation_participants VALUES (12, 8589938689, 12);
INSERT INTO public.history_operation_participants VALUES (13, 8589938689, 9);
INSERT INTO public.history_operation_participants VALUES (14, 8589942785, 12);
INSERT INTO public.history_operation_participants VALUES (15, 8589942785, 13);
INSERT INTO public.history_operation_participants VALUES (16, 8589946881, 12);
INSERT INTO public.history_operation_participants VALUES (17, 8589946881, 8);
INSERT INTO public.history_operation_participants VALUES (18, 8589950977, 12);
INSERT INTO public.history_operation_participants VALUES (19, 8589950977, 6);
INSERT INTO public.history_operation_participants VALUES (20, 8589955073, 10);
INSERT INTO public.history_operation_participants VALUES (21, 8589955073, 12);
INSERT INTO public.history_operation_participants VALUES (22, 8589959169, 7);
INSERT INTO public.history_operation_participants VALUES (23, 8589959169, 12);
INSERT INTO public.history_operation_participants VALUES (24, 8589963265, 12);
INSERT INTO public.history_operation_participants VALUES (25, 8589963265, 11);
INSERT INTO public.history_operation_participants VALUES (26, 8589967361, 12);
INSERT INTO public.history_operation_participants VALUES (27, 8589967361, 5);
INSERT INTO public.history_operation_participants VALUES (28, 8589971457, 12);
INSERT INTO public.history_operation_participants VALUES (29, 8589971457, 3);
INSERT INTO public.history_operation_participants VALUES (30, 8589975553, 12);
INSERT INTO public.history_operation_participants VALUES (31, 8589975553, 4);
INSERT INTO public.history_operation_participants VALUES (32, 8589979649, 12);
INSERT INTO public.history_operation_participants VALUES (33, 8589979649, 2);
INSERT INTO public.history_operation_participants VALUES (34, 8589983745, 12);
INSERT INTO public.history_operation_participants VALUES (35, 8589983745, 1);


--
-- Data for Name: history_operations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_operations VALUES (21474840577, 21474840576, 1, 5, '{"signer_key": "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN", "signer_weight": 1}', 'GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT');
INSERT INTO public.history_operations VALUES (21474844673, 21474844672, 1, 5, '{"signer_key": "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT", "signer_weight": 1}', 'GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN');
INSERT INTO public.history_operations VALUES (17179873281, 17179873280, 1, 5, '{"signer_key": "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ", "signer_weight": 1}', 'GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2');
INSERT INTO public.history_operations VALUES (17179877377, 17179877376, 1, 5, '{"signer_key": "GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2", "signer_weight": 1}', 'GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL');
INSERT INTO public.history_operations VALUES (17179881473, 17179881472, 1, 5, '{"signer_key": "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL", "signer_weight": 1}', 'GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ');
INSERT INTO public.history_operations VALUES (12884905985, 12884905984, 1, 5, '{"signer_key": "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y", "signer_weight": 1}', 'GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7');
INSERT INTO public.history_operations VALUES (12884910081, 12884910080, 1, 5, '{"signer_key": "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6", "signer_weight": 1}', 'GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ');
INSERT INTO public.history_operations VALUES (12884914177, 12884914176, 1, 1, '{"to": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "from": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "amount": "500.00000", "asset_type": "native"}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO public.history_operations VALUES (12884918273, 12884918272, 1, 5, '{"signer_key": "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ", "signer_weight": 1}', 'GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y');
INSERT INTO public.history_operations VALUES (12884922369, 12884922368, 1, 5, '{"signer_key": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "signer_weight": 1}', 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU');
INSERT INTO public.history_operations VALUES (8589938689, 8589938688, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU", "starting_balance": "10000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589942785, 8589942784, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2", "starting_balance": "10000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589946881, 8589946880, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBXGQJWVLWOYHFLVTKWV5FGHA3LNYY2JQKM7OAJAUEQFU6LPCSEFVXON", "starting_balance": "10000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589950977, 8589950976, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7", "starting_balance": "10000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589955073, 8589955072, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y", "starting_balance": "20000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589959169, 8589959168, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ", "starting_balance": "30000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589963265, 8589963264, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBS2EWY4PAVG4SOSGZ6CM5F4SIJUAIEXGKXEXCGTMCAK6FFVDDDXVHI6", "starting_balance": "40000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589967361, 8589967360, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ", "starting_balance": "50000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589971457, 8589971456, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2", "starting_balance": "60000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589975553, 8589975552, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL", "starting_balance": "70000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589979649, 8589979648, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN", "starting_balance": "80000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');
INSERT INTO public.history_operations VALUES (8589983745, 8589983744, 1, 0, '{"funder": "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H", "account": "GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT", "starting_balance": "90000.00000"}', 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H');


--
-- Data for Name: history_trades; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: history_transaction_participants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_transaction_participants VALUES (1, 21474840576, 1);
INSERT INTO public.history_transaction_participants VALUES (2, 21474844672, 2);
INSERT INTO public.history_transaction_participants VALUES (3, 17179873280, 3);
INSERT INTO public.history_transaction_participants VALUES (4, 17179877376, 4);
INSERT INTO public.history_transaction_participants VALUES (5, 17179881472, 5);
INSERT INTO public.history_transaction_participants VALUES (6, 12884905984, 6);
INSERT INTO public.history_transaction_participants VALUES (7, 12884910080, 7);
INSERT INTO public.history_transaction_participants VALUES (8, 12884914176, 9);
INSERT INTO public.history_transaction_participants VALUES (9, 12884914176, 8);
INSERT INTO public.history_transaction_participants VALUES (10, 12884918272, 10);
INSERT INTO public.history_transaction_participants VALUES (11, 12884922368, 9);
INSERT INTO public.history_transaction_participants VALUES (12, 8589938688, 12);
INSERT INTO public.history_transaction_participants VALUES (13, 8589938688, 9);
INSERT INTO public.history_transaction_participants VALUES (14, 8589942784, 12);
INSERT INTO public.history_transaction_participants VALUES (15, 8589942784, 13);
INSERT INTO public.history_transaction_participants VALUES (16, 8589946880, 12);
INSERT INTO public.history_transaction_participants VALUES (17, 8589946880, 8);
INSERT INTO public.history_transaction_participants VALUES (18, 8589950976, 12);
INSERT INTO public.history_transaction_participants VALUES (19, 8589950976, 6);
INSERT INTO public.history_transaction_participants VALUES (20, 8589955072, 12);
INSERT INTO public.history_transaction_participants VALUES (21, 8589955072, 10);
INSERT INTO public.history_transaction_participants VALUES (22, 8589959168, 12);
INSERT INTO public.history_transaction_participants VALUES (23, 8589959168, 7);
INSERT INTO public.history_transaction_participants VALUES (24, 8589963264, 11);
INSERT INTO public.history_transaction_participants VALUES (25, 8589963264, 12);
INSERT INTO public.history_transaction_participants VALUES (26, 8589967360, 12);
INSERT INTO public.history_transaction_participants VALUES (27, 8589967360, 5);
INSERT INTO public.history_transaction_participants VALUES (28, 8589971456, 12);
INSERT INTO public.history_transaction_participants VALUES (29, 8589971456, 3);
INSERT INTO public.history_transaction_participants VALUES (30, 8589975552, 12);
INSERT INTO public.history_transaction_participants VALUES (31, 8589975552, 4);
INSERT INTO public.history_transaction_participants VALUES (32, 8589979648, 12);
INSERT INTO public.history_transaction_participants VALUES (33, 8589979648, 2);
INSERT INTO public.history_transaction_participants VALUES (34, 8589983744, 12);
INSERT INTO public.history_transaction_participants VALUES (35, 8589983744, 1);


--
-- Data for Name: history_transactions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.history_transactions VALUES ('1cc851e62d2bdd34b6f0052754d8ec3e9fa1c3fbfdfabebb6c87c6ec682b3354', 5, 1, 'GCXOZEDULG6VUL7RMUPZDAMFAZ5WGXZDYNOGYLFK3YAPH2PIYGHSJXPT', 8589934593, 100, 1, '2019-06-13 14:07:46.503944', '2019-06-13 14:07:46.503969', 21474840576, 'AAAAAK7skHRZvVov8WUfkYGFBntjXyPDXGwsqt4A8+nowY8kAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAACGjfAcxwi5nMhbSC3prgEQlZYpKxP+vKkFf1uWngvhXAAAAAEAAAAAAAAAAejBjyQAAABAAGnrd1AuVMFbamNYkeg7fa+UT+YobWJWoIhyvXFSHaWpBfSQv52jBaXXJdHpGdeau/R7l2byN0AT/+V4oK94Bw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAUAAAAAAAAAAK7skHRZvVov8WUfkYGFBntjXyPDXGwsqt4A8+nowY8kAAAAAhhxGZwAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAAK7skHRZvVov8WUfkYGFBntjXyPDXGwsqt4A8+nowY8kAAAAAhhxGZwAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAACGjfAcxwi5nMhbSC3prgEQlZYpKxP+vKkFf1uWngvhXAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACu7JB0Wb1aL/FlH5GBhQZ7Y18jw1xsLKreAPPp6MGPJAAAAAIYcRoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACu7JB0Wb1aL/FlH5GBhQZ7Y18jw1xsLKreAPPp6MGPJAAAAAIYcRmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{AGnrd1AuVMFbamNYkeg7fa+UT+YobWJWoIhyvXFSHaWpBfSQv52jBaXXJdHpGdeau/R7l2byN0AT/+V4oK94Bw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('8c5d93e5b16da3f9ff541fbde785f619293f89648037805fcbfc978b98c52ec3', 5, 2, 'GCDI34A4Y4ELTHGILNEC32NOAEIJLFRJFMJ75PFJAV7VXFU6BPQVYTEN', 8589934593, 100, 1, '2019-06-13 14:07:46.505506', '2019-06-13 14:07:46.505521', 21474844672, 'AAAAAIaN8BzHCLmcyFtILemuARCVlikrE/68qQV/W5aeC+FcAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAACu7JB0Wb1aL/FlH5GBhQZ7Y18jw1xsLKreAPPp6MGPJAAAAAEAAAAAAAAAAZ4L4VwAAABAQq8VDt8safAiqJPFkp6rKDFSTLzZUERSSFxcx8LesiOA+Kf2XbGYy1SMYI+2CO/j7hEK6zZqKkbeCGj/AttfCw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAUAAAAAAAAAAIaN8BzHCLmcyFtILemuARCVlikrE/68qQV/W5aeC+FcAAAAAdzWT5wAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAUAAAAAAAAAAIaN8BzHCLmcyFtILemuARCVlikrE/68qQV/W5aeC+FcAAAAAdzWT5wAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAACu7JB0Wb1aL/FlH5GBhQZ7Y18jw1xsLKreAPPp6MGPJAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACGjfAcxwi5nMhbSC3prgEQlZYpKxP+vKkFf1uWngvhXAAAAAHc1lAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAFAAAAAAAAAACGjfAcxwi5nMhbSC3prgEQlZYpKxP+vKkFf1uWngvhXAAAAAHc1k+cAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Qq8VDt8safAiqJPFkp6rKDFSTLzZUERSSFxcx8LesiOA+Kf2XbGYy1SMYI+2CO/j7hEK6zZqKkbeCGj/AttfCw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('529434fdd8ba1a9aac4e3a241af5c1af59a560e144767af9ba95b20d89f6f0f9', 4, 1, 'GCBJRUEBARNP5HLASDN3ZNVQSYFFO2ETXCLZVCI2DSHTFBMA6R37U6S2', 8589934593, 100, 1, '2019-06-13 14:07:46.515246', '2019-06-13 14:07:46.515261', 17179873280, 'AAAAAIKY0IEEWv6dYJDbvLawlgpXaJO4l5qJGhyPMoWA9Hf6AAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAABU0ZA4qplYxjv/lGr/92Ewpx1ZhFvZkcpOE5Rupw1/hQAAAAEAAAAAAAAAAYD0d/oAAABA1eDMpqmrq3ILs9gvH45mUo446VRJ25ZKvnagrzTmXtEfoFF///rTBZ45LF/6cQonsBIR7EU2TOYDiCTUdYoYDg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAQAAAAAAAAAAIKY0IEEWv6dYJDbvLawlgpXaJO4l5qJGhyPMoWA9Hf6AAAAAWWgu5wAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAQAAAAAAAAAAIKY0IEEWv6dYJDbvLawlgpXaJO4l5qJGhyPMoWA9Hf6AAAAAWWgu5wAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAABU0ZA4qplYxjv/lGr/92Ewpx1ZhFvZkcpOE5Rupw1/hQAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACCmNCBBFr+nWCQ27y2sJYKV2iTuJeaiRocjzKFgPR3+gAAAAFloLwAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAACCmNCBBFr+nWCQ27y2sJYKV2iTuJeaiRocjzKFgPR3+gAAAAFloLucAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{1eDMpqmrq3ILs9gvH45mUo446VRJ25ZKvnagrzTmXtEfoFF///rTBZ45LF/6cQonsBIR7EU2TOYDiCTUdYoYDg==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('84032e6ec877b84de35686c65cb43304d6da2029891345a13332ba59dfbeeb06', 4, 2, 'GDSLCGMN4WK2SAOANYYOSIATOT2CWTSM5AHBQYAAXXJYEHI5CDEHRYIL', 8589934593, 100, 1, '2019-06-13 14:07:46.515396', '2019-06-13 14:07:46.51541', 17179877376, 'AAAAAOSxGY3llakBwG4w6SATdPQrTkzoDhhgAL3Tgh0dEMh4AAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAACCmNCBBFr+nWCQ27y2sJYKV2iTuJeaiRocjzKFgPR3+gAAAAEAAAAAAAAAAR0QyHgAAABARBapVp+ZYAhU04HBTJu92YtfDEXAwsYaGPGaY3s32XPxqZb/X2z+HkZFU+BYTlSq5v1eYxTkduCjq+BOT1GECg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAQAAAAAAAAAAOSxGY3llakBwG4w6SATdPQrTkzoDhhgAL3Tgh0dEMh4AAAAAaE7hZwAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAQAAAAAAAAAAOSxGY3llakBwG4w6SATdPQrTkzoDhhgAL3Tgh0dEMh4AAAAAaE7hZwAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAACCmNCBBFr+nWCQ27y2sJYKV2iTuJeaiRocjzKFgPR3+gAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAADksRmN5ZWpAcBuMOkgE3T0K05M6A4YYAC904IdHRDIeAAAAAGhO4YAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAADksRmN5ZWpAcBuMOkgE3T0K05M6A4YYAC904IdHRDIeAAAAAGhO4WcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{RBapVp+ZYAhU04HBTJu92YtfDEXAwsYaGPGaY3s32XPxqZb/X2z+HkZFU+BYTlSq5v1eYxTkduCjq+BOT1GECg==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('7014b5869a16bef76fe22b26fbff3d768904cd170748b0cd3f93aa6a7f8e3603', 4, 3, 'GBKNDEBYVKMVRRR376KGV77XMEYKOHKZQRN5TEOKJYJZI3VHBV7YKLJZ', 8589934593, 100, 1, '2019-06-13 14:07:46.515555', '2019-06-13 14:07:46.515569', 17179881472, 'AAAAAFTRkDiqmVjGO/+Uav/3YTCnHVmEW9mRyk4TlG6nDX+FAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAADksRmN5ZWpAcBuMOkgE3T0K05M6A4YYAC904IdHRDIeAAAAAEAAAAAAAAAAacNf4UAAABApJsBdLqYIsRmSpj920m41wvKcwQCWACdXYTM8NywQ2UaSqIiZj/agWX/nacbc/uQkhQwL+LLBVGELB4Ueu1HCw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAQAAAAAAAAAAFTRkDiqmVjGO/+Uav/3YTCnHVmEW9mRyk4TlG6nDX+FAAAAASoF8ZwAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAQAAAAAAAAAAFTRkDiqmVjGO/+Uav/3YTCnHVmEW9mRyk4TlG6nDX+FAAAAASoF8ZwAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAADksRmN5ZWpAcBuMOkgE3T0K05M6A4YYAC904IdHRDIeAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABU0ZA4qplYxjv/lGr/92Ewpx1ZhFvZkcpOE5Rupw1/hQAAAAEqBfIAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAEAAAAAAAAAABU0ZA4qplYxjv/lGr/92Ewpx1ZhFvZkcpOE5Rupw1/hQAAAAEqBfGcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{pJsBdLqYIsRmSpj920m41wvKcwQCWACdXYTM8NywQ2UaSqIiZj/agWX/nacbc/uQkhQwL+LLBVGELB4Ueu1HCw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('be88b2fa41df79d07fe66d987b34b2127198df0ea483dc80243f77072886dee3', 3, 1, 'GCJ6M7DRW5RZW73UEHEOGIGKLY5BWA5QNR6WX25G5BSYWDQYR77DSIH7', 8589934593, 100, 1, '2019-06-13 14:07:46.523742', '2019-06-13 14:07:46.523757', 12884905984, 'AAAAAJPmfHG3Y5t/dCHI4yDKXjobA7BsfWvrpuhliw4Yj/45AAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAADBmWo7Kaf/84fBBbk7Pa2oBNh2wjlgu5PrADV1Q4eDKAAAAAEAAAAAAAAAARiP/jkAAABA0/QhWG4CF9tL3qjbygLHFQ341sn9J/6Nutn562zWCtNJ0M8mt3oH2lQrVMHfZj04SFbgUTGxxN9IKp2Q+hzOAQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAMAAAAAAAAAAJPmfHG3Y5t/dCHI4yDKXjobA7BsfWvrpuhliw4Yj/45AAAAADuayZwAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAJPmfHG3Y5t/dCHI4yDKXjobA7BsfWvrpuhliw4Yj/45AAAAADuayZwAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAADBmWo7Kaf/84fBBbk7Pa2oBNh2wjlgu5PrADV1Q4eDKAAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAACT5nxxt2Obf3QhyOMgyl46GwOwbH1r66boZYsOGI/+OQAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACT5nxxt2Obf3QhyOMgyl46GwOwbH1r66boZYsOGI/+OQAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{0/QhWG4CF9tL3qjbygLHFQ341sn9J/6Nutn562zWCtNJ0M8mt3oH2lQrVMHfZj04SFbgUTGxxN9IKp2Q+hzOAQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('e4d3e43ac42240a68dcb9f93ad4f9ec20d6a1f6e633c5e93786f9b3f390d6e5e', 3, 2, 'GBSD7E6QRGGEBFWVF4VNSAWPQJW7FFLYEV7O3DZUIRLYATFMLCRJGSNQ', 8589934593, 100, 1, '2019-06-13 14:07:46.5239', '2019-06-13 14:07:46.523914', 12884910080, 'AAAAAGQ/k9CJjECW1S8q2QLPgm3ylXglfu2PNERXgEysWKKTAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAABlolsceCpuSdI2fCZ0vJITQCCXMq5LiNNggK8UtRjHegAAAAEAAAAAAAAAAaxYopMAAABAGLoP3FhgZ2HOOug5RjisdT/NL0rxXxc+gXmAG31D8ctpJ8W5uNL9QYMg+1CGn9IiQOJj+pvHnfUqORDVeqbSDg==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAMAAAAAAAAAAGQ/k9CJjECW1S8q2QLPgm3ylXglfu2PNERXgEysWKKTAAAAALLQXZwAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAGQ/k9CJjECW1S8q2QLPgm3ylXglfu2PNERXgEysWKKTAAAAALLQXZwAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAABlolsceCpuSdI2fCZ0vJITQCCXMq5LiNNggK8UtRjHegAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAABkP5PQiYxAltUvKtkCz4Jt8pV4JX7tjzREV4BMrFiikwAAAACy0F4AAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAABkP5PQiYxAltUvKtkCz4Jt8pV4JX7tjzREV4BMrFiikwAAAACy0F2cAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{GLoP3FhgZ2HOOug5RjisdT/NL0rxXxc+gXmAG31D8ctpJ8W5uNL9QYMg+1CGn9IiQOJj+pvHnfUqORDVeqbSDg==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('cebb875a00ff6e1383aef0fd251a76f22c1f9ab2a2dffcb077855736ade2659a', 3, 3, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934593, 100, 1, '2019-06-13 14:07:46.524043', '2019-06-13 14:07:46.524058', 12884914176, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAAAAAAAL68IAAAAAAAAAAAa7kvkwAAABA9Pu9pjykcRS60lqOLqN8FHz244QP8baYNeTTJZIlr3SbRC13qEr9uP4ORDgyCB/gcug2GKrDMuK0ST3QOaKUBw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAA=', 'AAAAAAAAAAEAAAAEAAAAAwAAAAIAAAAAAAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAADuaygAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAAD6VuoAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAMAAAAAAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAADuayTgAAAACAAAAAgAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAADif2LgAAAACAAAAAgAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msoAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{9Pu9pjykcRS60lqOLqN8FHz244QP8baYNeTTJZIlr3SbRC13qEr9uP4ORDgyCB/gcug2GKrDMuK0ST3QOaKUBw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('e0ea8d124661d812e2e209d572eae3e525e75bc20b16e6c277dbdfd0bf952c0e', 3, 4, 'GDAZS2R3FGT7744HYEC3SOZ5VWUAJWDWYI4WBO4T5MADK5KDQ6BSQZ4Y', 8589934593, 100, 1, '2019-06-13 14:07:46.524209', '2019-06-13 14:07:46.524223', 12884918272, 'AAAAAMGZajspp//zh8EFuTs9ragE2HbCOWC7k+sANXVDh4MoAAAAZAAAAAIAAAABAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAABkP5PQiYxAltUvKtkCz4Jt8pV4JX7tjzREV4BMrFiikwAAAAEAAAAAAAAAAUOHgygAAABAVnu2WoK32JLbRsWGUZsXaj8LgfsEJnLMtLw5YhUQ/Mpgt5uAXw+mj69dx1ZCluq8Ie63P5WHoQ+pCm13Xgq3Cw==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAMAAAAAAAAAAMGZajspp//zh8EFuTs9ragE2HbCOWC7k+sANXVDh4MoAAAAAHc1k5wAAAACAAAAAQAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAMGZajspp//zh8EFuTs9ragE2HbCOWC7k+sANXVDh4MoAAAAAHc1k5wAAAACAAAAAQAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAABkP5PQiYxAltUvKtkCz4Jt8pV4JX7tjzREV4BMrFiikwAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAACAAAAAAAAAADBmWo7Kaf/84fBBbk7Pa2oBNh2wjlgu5PrADV1Q4eDKAAAAAB3NZQAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAADBmWo7Kaf/84fBBbk7Pa2oBNh2wjlgu5PrADV1Q4eDKAAAAAB3NZOcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{Vnu2WoK32JLbRsWGUZsXaj8LgfsEJnLMtLw5YhUQ/Mpgt5uAXw+mj69dx1ZCluq8Ie63P5WHoQ+pCm13Xgq3Cw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('45b3dbce05efd82a906e0b7e9a4465fc5692542eecf861d86f6075410c639a3c', 3, 5, 'GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU', 8589934594, 100, 1, '2019-06-13 14:07:46.5244', '2019-06-13 14:07:46.524414', 12884922368, 'AAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAZAAAAAIAAAACAAAAAAAAAAAAAAABAAAAAAAAAAUAAAAAAAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAEAAAAAAAAAAa7kvkwAAABAuZRjg6HzR6UTgJukik8hMJD4b/LsH74hhf3Wqv893Ft8sjEoZQFbxeKOakiO+IlqLjTiCCu15UsgR3wpPCkWDQ==', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAFAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAwAAAAMAAAAAAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAADif2LgAAAACAAAAAgAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAADif2LgAAAACAAAAAgAAAAEAAAAAAAAAAAAAAAABAAAAAAAAAQAAAABuaCbVXZ2DlXWarV6UxwbW3GNJgpn3ASChIFp5bxSIWgAAAAEAAAAAAAAAAA==', 'AAAAAgAAAAMAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msmcAAAAAgAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAADAAAAAAAAAACuo3ot45qCPExpQ/3oHN+z17Ryis1lfMFYmQWgruS+TAAAAAA7msk4AAAAAgAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{uZRjg6HzR6UTgJukik8hMJD4b/LsH74hhf3Wqv893Ft8sjEoZQFbxeKOakiO+IlqLjTiCCu15UsgR3wpPCkWDQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d', 2, 1, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 1, 100, 1, '2019-06-13 14:07:46.536582', '2019-06-13 14:07:46.536597', 8589938688, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAABAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAArqN6LeOagjxMaUP96Bzfs9e0corNZXzBWJkFoK7kvkwAAAAAO5rKAAAAAAAAAAABVvwF9wAAAECDzqvkQBQoNAJifPRXDoLhvtycT3lFPCQ51gkdsFHaBNWw05S/VhW0Xgkr0CBPE4NaFV2Kmcs3ZwLmib4TRrML', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAK6jei3jmoI8TGlD/egc37PXtHKKzWV8wViZBaCu5L5MAAAAADuaygAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2s2vJMVAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAABAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{g86r5EAUKDQCYnz0Vw6C4b7cnE95RTwkOdYJHbBR2gTVsNOUv1YVtF4JK9AgTxODWhVdipnLN2cC5om+E0azCw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('164a5064eba64f2cdbadb856bf3448485fc626247ada3ed39cddf0f6902133b6', 2, 2, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 2, 100, 1, '2019-06-13 14:07:46.536793', '2019-06-13 14:07:46.536808', 8589942784, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAACAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAO2C/AO45YBD3tHVFO1R3A0MekP8JR6nN1A9eWidyItUAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEASEZiZbeFwCsrKBnKIus/05VtJDBrgosuhLQ/U6XUj4twWyhs7UtS4CMexOM6JqcfqJK10WlBkkwn4g8PIfjIG', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAADtgvwDuOWAQ97R1RTtUdwNDHpD/CUepzdQPXlonciLVAAAAADuaygAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2szAuZ1AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/+cAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{EhGYmW3hcArKygZyiLrP9OVbSQwa4KLLoS0P1Ol1I+LcFsobO1LUuAjHsTjOianH6iStdFpQZJMJ+IPDyH4yBg==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('2b2e82dbabb024b27a0c3140ca71d8ac9bc71831f9f5a3bd69eca3d88fb0ec5c', 2, 3, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 3, 100, 1, '2019-06-13 14:07:46.537037', '2019-06-13 14:07:46.537052', 8589946880, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAADAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAbmgm1V2dg5V1mq1elMcG1txjSYKZ9wEgoSBaeW8UiFoAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEDJul1tLGLF4Vxwt0dDCVEf6tb5l4byMrGgCp+lVZMmxct54iNf2mxtjx6Md5ZJ4E4Dlcsf46EAhBGSUPsn8fYD', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAG5oJtVdnYOVdZqtXpTHBtbcY0mCmfcBIKEgWnlvFIhaAAAAADuaygAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2svSTnVAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/84AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{ybpdbSxixeFccLdHQwlRH+rW+ZeG8jKxoAqfpVWTJsXLeeIjX9psbY8ejHeWSeBOA5XLH+OhAIQRklD7J/H2Aw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('dd06d15a6b7d9903d6dc1305c4209ebc78cdc93cd4b9059e0e4220ca814f7b03', 2, 4, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 4, 100, 1, '2019-06-13 14:07:46.537486', '2019-06-13 14:07:46.537501', 8589950976, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAEAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAk+Z8cbdjm390IcjjIMpeOhsDsGx9a+um6GWLDhiP/jkAAAAAO5rKAAAAAAAAAAABVvwF9wAAAEBHI1309PoxY1InoQX98iTdsNUbMt6wMm/DWo9AKUlKAMAo3dDUdcPW1Awrrd105Vl3Z2zPh92Xsm58k+gXC04G', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAJPmfHG3Y5t/dCHI4yDKXjobA7BsfWvrpuhliw4Yj/45AAAAADuaygAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2srj401AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/7UAAAAAAAAAAMAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/5wAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{RyNd9PT6MWNSJ6EF/fIk3bDVGzLesDJvw1qPQClJSgDAKN3Q1HXD1tQMK63ddOVZd2dsz4fdl7JufJPoFwtOBg==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('aeb2fa0145ccb6d6c0116aaa4a5f9a572e1e322de48bde5a636b0a0508ec34aa', 2, 5, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 5, 100, 1, '2019-06-13 14:07:46.53777', '2019-06-13 14:07:46.537785', 8589955072, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAFAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAwZlqOymn//OHwQW5Oz2tqATYdsI5YLuT6wA1dUOHgygAAAAAdzWUAAAAAAAAAAABVvwF9wAAAECbCCli2tJyShWmUD7M841/TETqa8LS4rjtbMWmnKLImyUysu3dpUty91J1LlDqoznVeY3HzDuOVe8BRUQvA5wJ', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAMGZajspp//zh8EFuTs9ragE2HbCOWC7k+sANXVDh4MoAAAAAHc1lAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2skHDP1AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/5wAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/4MAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{mwgpYtrSckoVplA+zPONf0xE6mvC0uK47WzFppyiyJslMrLt3aVLcvdSdS5Q6qM51XmNx8w7jlXvAUVELwOcCQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('f8a1dc5bd4939ff96fbc980e0e2c230cb9e7348b5db1855bfba4dd9e94c74903', 2, 6, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 6, 100, 1, '2019-06-13 14:07:46.538036', '2019-06-13 14:07:46.53805', 8589959168, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAGAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAZD+T0ImMQJbVLyrZAs+CbfKVeCV+7Y80RFeATKxYopMAAAAAstBeAAAAAAAAAAABVvwF9wAAAECgVKAxZlb+PiWsn0khQF6IU7Yy/tfGEjZR9iX1BxBrOo3yLP05xNoM7fylpRmT+VqsmkaRJEjawBQC+mO+YFcE', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAGQ/k9CJjECW1S8q2QLPgm3ylXglfu2PNERXgEysWKKTAAAAALLQXgAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2sY7y4VAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/4MAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/2oAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{oFSgMWZW/j4lrJ9JIUBeiFO2Mv7XxhI2UfYl9QcQazqN8iz9OcTaDO38paUZk/larJpGkSRI2sAUAvpjvmBXBA==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('976bc5c2c28a63765d7420a3492d4e03dcd9cebddeb1fd43718f17073b1323d5', 2, 7, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 7, 100, 1, '2019-06-13 14:07:46.538192', '2019-06-13 14:07:46.538206', 8589963264, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAHAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAZaJbHHgqbknSNnwmdLySE0AglzKuS4jTYICvFLUYx3oAAAAA7msoAAAAAAAAAAABVvwF9wAAAEB0VL3vmb6DftU3cVGm2sWk2SOYqLyVHxmxUxGifXYVW5OOsjTu6+mdJv0do7HZlIFZ0hKRP/O5wvLudOcs3/oN', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAGWiWxx4Km5J0jZ8JnS8khNAIJcyrkuI02CArxS1GMd6AAAAAO5rKAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2sKCHuVAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/2oAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/1EAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{dFS975m+g37VN3FRptrFpNkjmKi8lR8ZsVMRon12FVuTjrI07uvpnSb9HaOx2ZSBWdISkT/zucLy7nTnLN/6DQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('4bd3123606380202f65f9fce531ab6ce03461ffcf491a85a9b7abcaba01d3691', 2, 8, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 8, 100, 1, '2019-06-13 14:07:46.538347', '2019-06-13 14:07:46.538361', 8589967360, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAIAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAVNGQOKqZWMY7/5Rq//dhMKcdWYRb2ZHKThOUbqcNf4UAAAABKgXyAAAAAAAAAAABVvwF9wAAAEALWpdAp52+0A0d0PtC8dKiJyaQEHOpVSmKjMFPl6zDQgxepiYoaZncFdXreBjztgCOHKfQNBiyC2RZz8mPnTgN', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAFTRkDiqmVjGO/+Uav/3YTCnHVmEW9mRyk4TlG6nDX+FAAAAASoF8gAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2r3aBx1AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/1EAAAAAAAAAAcAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/zgAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{C1qXQKedvtANHdD7QvHSoicmkBBzqVUpiozBT5esw0IMXqYmKGmZ3BXV63gY87YAjhyn0DQYsgtkWc/Jj504DQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('72ec601c273c9f1a6ee484885e4412a1cbc5d0968030a641921c982edcb23ebe', 2, 9, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 9, 100, 1, '2019-06-13 14:07:46.538655', '2019-06-13 14:07:46.53867', 8589971456, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAJAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAgpjQgQRa/p1gkNu8trCWCldok7iXmokaHI8yhYD0d/oAAAABZaC8AAAAAAAAAAABVvwF9wAAAECp4/VeLFiiUBziqo02xmv7TZwJmq0Ik1hOn9RNeUyCT7EfvH5Ubs9uIKlpFsc1+2blOatby9d/8x1JUJUUJxoJ', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAIKY0IEEWv6dYJDbvLawlgpXaJO4l5qJGhyPMoWA9Hf6AAAAAWWgvAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2rhDhC1AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/zgAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/x8AAAAAAAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{qeP1XixYolAc4qqNNsZr+02cCZqtCJNYTp/UTXlMgk+xH7x+VG7PbiCpaRbHNftm5TmrW8vXf/MdSVCVFCcaCQ==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('4c3f4f97742bbecae73a6243a171a2fca12d61a30efa052f34d19d7fd3403877', 2, 10, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 10, 100, 1, '2019-06-13 14:07:46.538858', '2019-06-13 14:07:46.538874', 8589975552, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAKAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAA5LEZjeWVqQHAbjDpIBN09CtOTOgOGGAAvdOCHR0QyHgAAAABoTuGAAAAAAAAAAABVvwF9wAAAECwZuhsAmu9Mqnf4EyMFkaB20nnUsI5ViXZUWNFQuLbT8PWFJXG5pPK/yn2dGZYfZzROe4rAI/A/Ro80/zYPXwC', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAOSxGY3llakBwG4w6SATdPQrTkzoDhhgAL3Tgh0dEMh4AAAAAaE7hgAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2rG+lhVAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/x8AAAAAAAAAAkAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/wYAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{sGbobAJrvTKp3+BMjBZGgdtJ51LCOVYl2VFjRULi20/D1hSVxuaTyv8p9nRmWH2c0TnuKwCPwP0aPNP82D18Ag==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('3a64027bbf29d4f65944d9b0c1603a0ff046445d936fd300b8ecab1902c39bd2', 2, 11, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 11, 100, 1, '2019-06-13 14:07:46.539151', '2019-06-13 14:07:46.539166', 8589979648, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAALAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAho3wHMcIuZzIW0gt6a4BEJWWKSsT/rypBX9blp4L4VwAAAAB3NZQAAAAAAAAAAABVvwF9wAAAECWaOr2TdjnBbdliFjX7NZBgP36spCjjviXwczFmzuJIBE4hT8ZD5mR5d9t8nHwC9kXZ/vahFxtqXQd9TnYtzoL', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAIaN8BzHCLmcyFtILemuARCVlikrE/68qQV/W5aeC+FcAAAAAdzWUAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2qpLPNVAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/wYAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/u0AAAAAAAAAAsAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{lmjq9k3Y5wW3ZYhY1+zWQYD9+rKQo474l8HMxZs7iSAROIU/GQ+ZkeXfbfJx8AvZF2f72oRcbal0HfU52Lc6Cw==}', 'none', NULL, NULL);
INSERT INTO public.history_transactions VALUES ('feab91c968afd32ba271fe60dc013a73a44b759723db15f4d06459185363097d', 2, 12, 'GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H', 12, 100, 1, '2019-06-13 14:07:46.539304', '2019-06-13 14:07:46.539319', 8589983744, 'AAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3AAAAZAAAAAAAAAAMAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAruyQdFm9Wi/xZR+RgYUGe2NfI8NcbCyq3gDz6ejBjyQAAAACGHEaAAAAAAAAAAABVvwF9wAAAEAbnGxhPfJAixWLOGhHe2sowXQrQd5d1sTmM+ydZDWZnD49Hx3cpubO7CH2VaDJaX8BTn94mkDlyAi2kc6dIFkE', 'AAAAAAAAAGQAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAA=', 'AAAAAAAAAAEAAAACAAAAAAAAAAIAAAAAAAAAAK7skHRZvVov8WUfkYGFBntjXyPDXGwsqt4A8+nowY8kAAAAAhhxGgAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAIAAAAAAAAAAGL8HQvQkbK2HA3WVjRrKmjX00fG8sLI7m0ERwJW/AX3DeC2qHpeG1AAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAA', 'AAAAAgAAAAMAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/u0AAAAAAAAAAsAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAEAAAACAAAAAAAAAABi/B0L0JGythwN1lY0aypo19NHxvLCyO5tBEcCVvwF9w3gtrOnY/tQAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==', '{G5xsYT3yQIsVizhoR3trKMF0K0HeXdbE5jPsnWQ1mZw+PR8d3Kbmzuwh9lWgyWl/AU5/eJpA5cgItpHOnSBZBA==}', 'none', NULL, NULL);


--
-- Name: history_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.history_accounts_id_seq', 13, true);


--
-- Name: history_assets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.history_assets_id_seq', 1, false);


--
-- Name: history_operation_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.history_operation_participants_id_seq', 35, true);


--
-- Name: history_transaction_participants_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.history_transaction_participants_id_seq', 35, true);


--
-- Name: asset_stats asset_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.asset_stats
    ADD CONSTRAINT asset_stats_pkey PRIMARY KEY (id);


--
-- Name: gorp_migrations gorp_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.gorp_migrations
    ADD CONSTRAINT gorp_migrations_pkey PRIMARY KEY (id);


--
-- Name: history_assets history_assets_asset_code_asset_type_asset_issuer_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_assets
    ADD CONSTRAINT history_assets_asset_code_asset_type_asset_issuer_key UNIQUE (asset_code, asset_type, asset_issuer);


--
-- Name: history_assets history_assets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_assets
    ADD CONSTRAINT history_assets_pkey PRIMARY KEY (id);


--
-- Name: history_operation_participants history_operation_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_operation_participants
    ADD CONSTRAINT history_operation_participants_pkey PRIMARY KEY (id);


--
-- Name: history_transaction_participants history_transaction_participants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_transaction_participants
    ADD CONSTRAINT history_transaction_participants_pkey PRIMARY KEY (id);


--
-- Name: asset_by_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX asset_by_code ON public.history_assets USING btree (asset_code);


--
-- Name: asset_by_issuer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX asset_by_issuer ON public.history_assets USING btree (asset_issuer);


--
-- Name: by_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_account ON public.history_transactions USING btree (account, account_sequence);


--
-- Name: by_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_hash ON public.history_transactions USING btree (transaction_hash);


--
-- Name: by_ledger; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX by_ledger ON public.history_transactions USING btree (ledger_sequence, application_order);


--
-- Name: hist_e_by_order; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_e_by_order ON public.history_effects USING btree (history_operation_id, "order");


--
-- Name: hist_e_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_e_id ON public.history_effects USING btree (history_account_id, history_operation_id, "order");


--
-- Name: hist_op_p_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_op_p_id ON public.history_operation_participants USING btree (history_account_id, history_operation_id);


--
-- Name: hist_tx_p_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hist_tx_p_id ON public.history_transaction_participants USING btree (history_account_id, history_transaction_id);


--
-- Name: hop_by_hoid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX hop_by_hoid ON public.history_operation_participants USING btree (history_operation_id);


--
-- Name: hs_ledger_by_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hs_ledger_by_id ON public.history_ledgers USING btree (id);


--
-- Name: hs_transaction_by_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX hs_transaction_by_id ON public.history_transactions USING btree (id);


--
-- Name: htp_by_htid; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htp_by_htid ON public.history_transaction_participants USING btree (history_transaction_id);


--
-- Name: htrd_by_base_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_base_account ON public.history_trades USING btree (base_account_id);


--
-- Name: htrd_by_base_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_base_offer ON public.history_trades USING btree (base_offer_id);


--
-- Name: htrd_by_counter_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_counter_account ON public.history_trades USING btree (counter_account_id);


--
-- Name: htrd_by_counter_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_counter_offer ON public.history_trades USING btree (counter_offer_id);


--
-- Name: htrd_by_offer; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_by_offer ON public.history_trades USING btree (offer_id);


--
-- Name: htrd_counter_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_counter_lookup ON public.history_trades USING btree (counter_asset_id);


--
-- Name: htrd_pair_time_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_pair_time_lookup ON public.history_trades USING btree (base_asset_id, counter_asset_id, ledger_closed_at);


--
-- Name: htrd_pid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX htrd_pid ON public.history_trades USING btree (history_operation_id, "order");


--
-- Name: htrd_time_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX htrd_time_lookup ON public.history_trades USING btree (ledger_closed_at);


--
-- Name: index_history_accounts_on_address; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_accounts_on_address ON public.history_accounts USING btree (address);


--
-- Name: index_history_accounts_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_accounts_on_id ON public.history_accounts USING btree (id);


--
-- Name: index_history_effects_on_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_effects_on_type ON public.history_effects USING btree (type);


--
-- Name: index_history_ledgers_on_closed_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_ledgers_on_closed_at ON public.history_ledgers USING btree (closed_at);


--
-- Name: index_history_ledgers_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_id ON public.history_ledgers USING btree (id);


--
-- Name: index_history_ledgers_on_importer_version; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_ledgers_on_importer_version ON public.history_ledgers USING btree (importer_version);


--
-- Name: index_history_ledgers_on_ledger_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_ledger_hash ON public.history_ledgers USING btree (ledger_hash);


--
-- Name: index_history_ledgers_on_previous_ledger_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_previous_ledger_hash ON public.history_ledgers USING btree (previous_ledger_hash);


--
-- Name: index_history_ledgers_on_sequence; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_ledgers_on_sequence ON public.history_ledgers USING btree (sequence);


--
-- Name: index_history_operations_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_operations_on_id ON public.history_operations USING btree (id);


--
-- Name: index_history_operations_on_transaction_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_operations_on_transaction_id ON public.history_operations USING btree (transaction_id);


--
-- Name: index_history_operations_on_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_history_operations_on_type ON public.history_operations USING btree (type);


--
-- Name: index_history_transactions_on_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_history_transactions_on_id ON public.history_transactions USING btree (id);


--
-- Name: trade_effects_by_order_book; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX trade_effects_by_order_book ON public.history_effects USING btree (((details ->> 'sold_asset_type'::text)), ((details ->> 'sold_asset_code'::text)), ((details ->> 'sold_asset_issuer'::text)), ((details ->> 'bought_asset_type'::text)), ((details ->> 'bought_asset_code'::text)), ((details ->> 'bought_asset_issuer'::text))) WHERE (type = 33);


--
-- Name: asset_stats asset_stats_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.asset_stats
    ADD CONSTRAINT asset_stats_id_fkey FOREIGN KEY (id) REFERENCES public.history_assets(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: history_trades history_trades_base_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_trades
    ADD CONSTRAINT history_trades_base_account_id_fkey FOREIGN KEY (base_account_id) REFERENCES public.history_accounts(id);


--
-- Name: history_trades history_trades_base_asset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_trades
    ADD CONSTRAINT history_trades_base_asset_id_fkey FOREIGN KEY (base_asset_id) REFERENCES public.history_assets(id);


--
-- Name: history_trades history_trades_counter_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_trades
    ADD CONSTRAINT history_trades_counter_account_id_fkey FOREIGN KEY (counter_account_id) REFERENCES public.history_accounts(id);


--
-- Name: history_trades history_trades_counter_asset_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.history_trades
    ADD CONSTRAINT history_trades_counter_asset_id_fkey FOREIGN KEY (counter_asset_id) REFERENCES public.history_assets(id);


--
-- PostgreSQL database dump complete
--

