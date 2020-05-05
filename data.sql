--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Debian 12.2-2.pgdg100+1)
-- Dumped by pg_dump version 12.2 (Debian 12.2-2.pgdg100+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.accounts (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    email text,
    password text,
    github_id text,
    company_name text,
    company_website text,
    account_token text
);


ALTER TABLE public.accounts OWNER TO postgres;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO postgres;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: apps; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.apps (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    account_id integer,
    app_name text,
    access_url text,
    registry_download_url text,
    local_access_url text,
    git_url text
);


ALTER TABLE public.apps OWNER TO postgres;

--
-- Name: apps_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.apps_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.apps_id_seq OWNER TO postgres;

--
-- Name: apps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.apps_id_seq OWNED BY public.apps.id;


--
-- Name: deployment_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deployment_settings (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    replicas integer,
    plan_id integer
);


ALTER TABLE public.deployment_settings OWNER TO postgres;

--
-- Name: deployment_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deployment_settings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.deployment_settings_id_seq OWNER TO postgres;

--
-- Name: deployment_settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deployment_settings_id_seq OWNED BY public.deployment_settings.id;


--
-- Name: domains; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.domains (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    address text
);


ALTER TABLE public.domains OWNER TO postgres;

--
-- Name: domains_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.domains_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.domains_id_seq OWNER TO postgres;

--
-- Name: domains_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.domains_id_seq OWNED BY public.domains.id;


--
-- Name: environments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.environments (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    res_id integer,
    env_key text,
    env_value text
);


ALTER TABLE public.environments OWNER TO postgres;

--
-- Name: environments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.environments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.environments_id_seq OWNER TO postgres;

--
-- Name: environments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.environments_id_seq OWNED BY public.environments.id;


--
-- Name: plans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plans (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    name text,
    alias text,
    price numeric
);


ALTER TABLE public.plans OWNER TO postgres;

--
-- Name: plans_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.plans_id_seq OWNER TO postgres;

--
-- Name: plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.plans_id_seq OWNED BY public.plans.id;


--
-- Name: quota; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.quota (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    memory numeric,
    cpu numeric,
    storage_size numeric
);


ALTER TABLE public.quota OWNER TO postgres;

--
-- Name: quota_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.quota_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.quota_id_seq OWNER TO postgres;

--
-- Name: quota_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.quota_id_seq OWNED BY public.quota.id;


--
-- Name: releases; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.releases (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    last_check_sum text,
    version_number bigint,
    docker_url text
);


ALTER TABLE public.releases OWNER TO postgres;

--
-- Name: releases_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.releases_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.releases_id_seq OWNER TO postgres;

--
-- Name: releases_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.releases_id_seq OWNED BY public.releases.id;


--
-- Name: resource_configs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_configs (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    resource_id integer
);


ALTER TABLE public.resource_configs OWNER TO postgres;

--
-- Name: resource_configs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_configs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_configs_id_seq OWNER TO postgres;

--
-- Name: resource_configs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resource_configs_id_seq OWNED BY public.resource_configs.id;


--
-- Name: resource_envs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resource_envs (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    resource_id integer,
    env_key text,
    env_value text
);


ALTER TABLE public.resource_envs OWNER TO postgres;

--
-- Name: resource_envs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resource_envs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resource_envs_id_seq OWNER TO postgres;

--
-- Name: resource_envs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resource_envs_id_seq OWNED BY public.resource_envs.id;


--
-- Name: resources; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.resources (
    id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    app_id integer,
    name text
);


ALTER TABLE public.resources OWNER TO postgres;

--
-- Name: resources_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.resources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resources_id_seq OWNER TO postgres;

--
-- Name: resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.resources_id_seq OWNED BY public.resources.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: apps id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apps ALTER COLUMN id SET DEFAULT nextval('public.apps_id_seq'::regclass);


--
-- Name: deployment_settings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_settings ALTER COLUMN id SET DEFAULT nextval('public.deployment_settings_id_seq'::regclass);


--
-- Name: domains id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.domains ALTER COLUMN id SET DEFAULT nextval('public.domains_id_seq'::regclass);


--
-- Name: environments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.environments ALTER COLUMN id SET DEFAULT nextval('public.environments_id_seq'::regclass);


--
-- Name: plans id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plans ALTER COLUMN id SET DEFAULT nextval('public.plans_id_seq'::regclass);


--
-- Name: quota id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quota ALTER COLUMN id SET DEFAULT nextval('public.quota_id_seq'::regclass);


--
-- Name: releases id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.releases ALTER COLUMN id SET DEFAULT nextval('public.releases_id_seq'::regclass);


--
-- Name: resource_configs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_configs ALTER COLUMN id SET DEFAULT nextval('public.resource_configs_id_seq'::regclass);


--
-- Name: resource_envs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_envs ALTER COLUMN id SET DEFAULT nextval('public.resource_envs_id_seq'::regclass);


--
-- Name: resources id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resources ALTER COLUMN id SET DEFAULT nextval('public.resources_id_seq'::regclass);


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.accounts (id, created_at, updated_at, deleted_at, name, email, password, github_id, company_name, company_website, account_token) FROM stdin;
1	2020-03-20 09:44:43.81531+00	2020-03-20 09:44:43.81531+00	\N	Lekan Adigun	adigunadunfe@gmail.com	$2a$10$IUSleMKgiT9kF4SqLUb.YOT4.0jw23TCYelSmhQIqDTRbi67uJryu				rQKUcjWRjkPNSVguipjihtoAgAmOTHpk:SldzySaVOEHqaKWTitoRUhgqLBOqOnUVNKuL
\.


--
-- Data for Name: apps; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.apps (id, created_at, updated_at, deleted_at, account_id, app_name, access_url, registry_download_url, local_access_url, git_url) FROM stdin;
3	2020-03-20 09:50:32.509332+00	2020-03-20 09:50:32.509332+00	\N	1	snowy-mountain	https://snowy-mountain.hostgoapp.com			https://git.hostgoapp.com/snowy-mountain.git
4	2020-04-17 12:56:36.634672+00	2020-04-17 12:56:36.634672+00	\N	1	broken-field	https://broken-field.hostgoapp.com			https://git.hostgoapp.com/broken-field.git
5	2020-04-17 13:04:08.028717+00	2020-04-17 13:04:08.028717+00	\N	1	bold-dew	https://bold-dew.hostgoapp.com			https://git.hostgoapp.com/bold-dew.git
6	2020-04-19 09:02:17.443879+00	2020-04-19 09:02:17.443879+00	\N	1	wild-cloud	https://wild-cloud.hostgoapp.com			https://git.hostgoapp.com/wild-cloud.git
7	2020-04-19 20:29:51.244512+00	2020-04-19 20:29:51.244512+00	\N	1	still-day-b312c208	https://still-day-b312c208.hostgoapp.com			https://git.hostgoapp.com/still-day-b312c208.git
8	2020-04-19 20:29:53.211141+00	2020-04-19 20:29:53.211141+00	\N	1	bitter-firefly-54d30848	https://bitter-firefly-54d30848.hostgoapp.com			https://git.hostgoapp.com/bitter-firefly-54d30848.git
9	2020-04-19 20:29:54.18607+00	2020-04-19 20:29:54.18607+00	\N	1	summer-month-1a3034a7	https://summer-month-1a3034a7.hostgoapp.com			https://git.hostgoapp.com/summer-month-1a3034a7.git
10	2020-04-19 20:29:55.015445+00	2020-04-19 20:29:55.015445+00	\N	1	ancient-cherry-36ee77a2	https://ancient-cherry-36ee77a2.hostgoapp.com			https://git.hostgoapp.com/ancient-cherry-36ee77a2.git
11	2020-04-19 20:29:55.809409+00	2020-04-19 20:29:55.809409+00	\N	1	muddy-shadow-129fcb2c	https://muddy-shadow-129fcb2c.hostgoapp.com			https://git.hostgoapp.com/muddy-shadow-129fcb2c.git
\.


--
-- Data for Name: deployment_settings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.deployment_settings (id, created_at, updated_at, deleted_at, app_id, replicas, plan_id) FROM stdin;
1	2020-03-20 09:59:50.965872+00	2020-03-20 09:59:50.965872+00	\N	3	1	0
2	2020-04-17 13:05:08.351489+00	2020-04-17 13:05:08.351489+00	\N	5	1	0
3	2020-04-19 09:02:29.995038+00	2020-04-19 09:02:29.995038+00	\N	6	1	0
4	2020-04-19 20:34:16.303184+00	2020-04-19 20:34:16.303184+00	\N	11	1	0
\.


--
-- Data for Name: domains; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.domains (id, created_at, updated_at, deleted_at, app_id, address) FROM stdin;
1	2020-04-17 13:16:48.852467+00	2020-04-17 13:16:48.852467+00	\N	5	note.local
\.


--
-- Data for Name: environments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.environments (id, created_at, updated_at, deleted_at, app_id, res_id, env_key, env_value) FROM stdin;
111	2020-04-19 19:27:30.607081+00	2020-04-19 19:27:30.607081+00	2020-04-19 19:33:08.663458+00	6	22	MONGO_INITDB_ROOT_PASSWORD	f0683cb808925faa147eb73d3a0f51b9
1	2020-03-20 10:05:32.537262+00	2020-03-20 10:05:32.537262+00	2020-03-22 09:18:05.43236+00	3	1	POSTGRES_USER	WseeVnEZfcKBYSDjbbQwpQSmpRjBoFfqiba
2	2020-03-20 10:05:32.54342+00	2020-03-20 10:05:32.54342+00	2020-03-22 09:18:05.43236+00	3	1	POSTGRES_PASSWORD	04114d7d567340021d3c3911883a46e9
3	2020-03-20 10:05:32.544989+00	2020-03-20 10:05:32.544989+00	2020-03-22 09:18:05.43236+00	3	1	POSTGRES_DB	WseeVnEZfcKBYSDjbbQw
4	2020-03-20 10:05:32.546068+00	2020-03-20 10:05:32.546068+00	2020-03-22 09:18:05.43236+00	3	1	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
5	2020-03-20 10:05:32.682492+00	2020-03-20 10:05:32.682492+00	2020-03-22 09:18:05.43236+00	3	1	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
6	2020-03-22 09:18:16.148929+00	2020-03-22 09:18:16.148929+00	2020-03-22 16:22:49.424739+00	3	2	POSTGRES_USER	yGqGTsQJHJnWjZYqWlWlpIcjQoTgZZrwUPT
7	2020-03-22 09:18:16.150961+00	2020-03-22 09:18:16.150961+00	2020-03-22 16:22:49.424739+00	3	2	POSTGRES_PASSWORD	46bfa6aecacbd4dc98179e3b92022edf
8	2020-03-22 09:18:16.152035+00	2020-03-22 09:18:16.152035+00	2020-03-22 16:22:49.424739+00	3	2	POSTGRES_DB	yGqGTsQJHJnWjZYqWlWl
9	2020-03-22 09:18:16.152714+00	2020-03-22 09:18:16.152714+00	2020-03-22 16:22:49.424739+00	3	2	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
10	2020-03-22 09:18:16.329189+00	2020-03-22 09:18:16.329189+00	2020-03-22 16:22:49.424739+00	3	2	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
11	2020-03-22 16:23:03.243766+00	2020-03-22 16:23:03.243766+00	2020-03-22 16:27:40.785455+00	3	3	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
12	2020-03-22 16:23:03.245907+00	2020-03-22 16:23:03.245907+00	2020-03-22 16:27:40.785455+00	3	3	POSTGRES_USER	ssUdVELtdWDxWACwXyHySiHMNDbxybalbCV
13	2020-03-22 16:23:03.246517+00	2020-03-22 16:23:03.246517+00	2020-03-22 16:27:40.785455+00	3	3	POSTGRES_PASSWORD	a1177b3192fa1a082d5d413182d57bc0
14	2020-03-22 16:23:03.247046+00	2020-03-22 16:23:03.247046+00	2020-03-22 16:27:40.785455+00	3	3	POSTGRES_DB	ssUdVELtdWDxWACwXyHy
15	2020-03-22 16:23:03.334104+00	2020-03-22 16:23:03.334104+00	2020-03-22 16:27:40.785455+00	3	3	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
16	2020-03-22 16:27:26.766915+00	2020-03-22 16:27:26.766915+00	2020-03-22 16:38:37.829655+00	3	0	KEY	VALUE
17	2020-03-22 16:27:53.568791+00	2020-03-22 16:27:53.568791+00	2020-03-22 16:38:59.25625+00	3	4	POSTGRES_DB	iwFNkJoFMTxGFLBdYgDA
18	2020-03-22 16:27:53.569211+00	2020-03-22 16:27:53.569211+00	2020-03-22 16:38:59.25625+00	3	4	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
19	2020-03-22 16:27:53.569626+00	2020-03-22 16:27:53.569626+00	2020-03-22 16:38:59.25625+00	3	4	POSTGRES_USER	iwFNkJoFMTxGFLBdYgDAjsqBuHNDatCKivS
20	2020-03-22 16:27:53.570067+00	2020-03-22 16:27:53.570067+00	2020-03-22 16:38:59.25625+00	3	4	POSTGRES_PASSWORD	7fb8254b90b992be437bfd0494f645e2
21	2020-03-22 16:27:53.613756+00	2020-03-22 16:27:53.613756+00	2020-03-22 16:38:59.25625+00	3	4	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
22	2020-03-22 16:39:25.44462+00	2020-03-22 16:39:25.44462+00	2020-03-22 16:56:05.338943+00	3	5	POSTGRES_USER	fkwHZjisSyBANBNGpaVEDkeZwnfxqaOMvlk
23	2020-03-22 16:39:25.44573+00	2020-03-22 16:39:25.44573+00	2020-03-22 16:56:05.338943+00	3	5	POSTGRES_PASSWORD	c699b81e65ef045f763cdba1ad5059ce
24	2020-03-22 16:39:25.446296+00	2020-03-22 16:39:25.446296+00	2020-03-22 16:56:05.338943+00	3	5	POSTGRES_DB	fkwHZjisSyBANBNGpaVE
25	2020-03-22 16:39:25.44691+00	2020-03-22 16:39:25.44691+00	2020-03-22 16:56:05.338943+00	3	5	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
26	2020-03-22 16:39:25.523046+00	2020-03-22 16:39:25.523046+00	2020-03-22 16:56:05.338943+00	3	5	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
28	2020-03-22 16:56:27.923491+00	2020-03-22 16:56:27.923491+00	2020-03-22 17:03:12.539353+00	3	6	POSTGRES_PASSWORD	fead6ee70fcb9fcfdc56d83078b0533a
29	2020-03-22 16:56:27.934277+00	2020-03-22 16:56:27.934277+00	2020-03-22 17:03:12.539353+00	3	6	POSTGRES_DB	paybTwtzIZmZDBkczCFK
30	2020-03-22 16:56:27.93517+00	2020-03-22 16:56:27.93517+00	2020-03-22 17:03:12.539353+00	3	6	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
31	2020-03-22 16:56:27.935806+00	2020-03-22 16:56:27.935806+00	2020-03-22 17:03:12.539353+00	3	6	POSTGRES_USER	paybTwtzIZmZDBkczCFKRRLdEbWsNnmHoHf
32	2020-03-22 16:56:28.055942+00	2020-03-22 16:56:28.055942+00	2020-03-22 17:03:12.539353+00	3	6	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
27	2020-03-22 16:39:35.580605+00	2020-03-22 16:39:35.580605+00	2020-03-22 17:03:41.618866+00	3	0	KEY	VALUE
38	2020-03-22 17:03:41.62126+00	2020-03-22 17:03:41.62126+00	\N	3	0	KEY	VALUE
39	2020-03-22 17:03:47.069582+00	2020-03-22 17:03:47.069582+00	\N	3	0	KEY1	VALUE1
33	2020-03-22 17:03:23.952837+00	2020-03-22 17:03:23.952837+00	2020-04-10 19:43:22.231149+00	3	7	POSTGRES_USER	XHFsHuDkrHUASxnMKIcXsVTGilbtNMVQJXo
34	2020-03-22 17:03:23.953332+00	2020-03-22 17:03:23.953332+00	2020-04-10 19:43:22.231149+00	3	7	POSTGRES_PASSWORD	5ae00524f4e14d529561f6f7bb3d1b90
35	2020-03-22 17:03:23.955084+00	2020-03-22 17:03:23.955084+00	2020-04-10 19:43:22.231149+00	3	7	POSTGRES_DB	XHFsHuDkrHUASxnMKIcX
36	2020-03-22 17:03:23.955588+00	2020-03-22 17:03:23.955588+00	2020-04-10 19:43:22.231149+00	3	7	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
37	2020-03-22 17:03:24.018635+00	2020-03-22 17:03:24.018635+00	2020-04-10 19:43:22.231149+00	3	7	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
40	2020-04-10 19:43:41.328716+00	2020-04-10 19:43:41.328716+00	2020-04-10 19:51:55.264153+00	3	8	POSTGRES_DB	tzpskISuaOAyeNjoiNQi
41	2020-04-10 19:43:41.33131+00	2020-04-10 19:43:41.33131+00	2020-04-10 19:51:55.264153+00	3	8	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
42	2020-04-10 19:43:41.331739+00	2020-04-10 19:43:41.331739+00	2020-04-10 19:51:55.264153+00	3	8	POSTGRES_USER	tzpskISuaOAyeNjoiNQiEBPRnRaQuMouTti
43	2020-04-10 19:43:41.332126+00	2020-04-10 19:43:41.332126+00	2020-04-10 19:51:55.264153+00	3	8	POSTGRES_PASSWORD	4a9f5a606307a21e5df34faf211370de
44	2020-04-10 19:52:04.168131+00	2020-04-10 19:52:04.168131+00	2020-04-10 19:57:13.531423+00	3	9	POSTGRES_USER	oshafHfyeceAdcfJNJtjZoQvuMwoDZnImIJ
45	2020-04-10 19:52:04.168519+00	2020-04-10 19:52:04.168519+00	2020-04-10 19:57:13.531423+00	3	9	POSTGRES_PASSWORD	e3e4d353527322352b170bef57dead4c
46	2020-04-10 19:52:04.168973+00	2020-04-10 19:52:04.168973+00	2020-04-10 19:57:13.531423+00	3	9	POSTGRES_DB	oshafHfyeceAdcfJNJtj
47	2020-04-10 19:52:04.169391+00	2020-04-10 19:52:04.169391+00	2020-04-10 19:57:13.531423+00	3	9	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
48	2020-04-10 19:57:18.12345+00	2020-04-10 19:57:18.12345+00	\N	3	10	POSTGRES_USER	RnSSczXJdhcBBcNcRAmzmgjIozyOYdphFgf
49	2020-04-10 19:57:18.132657+00	2020-04-10 19:57:18.132657+00	\N	3	10	POSTGRES_PASSWORD	d89732d4198ca1da454d7d2aaf07b562
50	2020-04-10 19:57:18.133571+00	2020-04-10 19:57:18.133571+00	\N	3	10	POSTGRES_DB	RnSSczXJdhcBBcNcRAmz
51	2020-04-10 19:57:18.134387+00	2020-04-10 19:57:18.134387+00	\N	3	10	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
52	2020-04-10 19:57:18.230317+00	2020-04-10 19:57:18.230317+00	\N	3	10	PG_HOST	svc-pg-snowy-mountain.namespace-storm:5432
53	2020-04-11 09:50:54.242123+00	2020-04-11 09:50:54.242123+00	2020-04-11 10:01:43.739232+00	3	11	MYSQL_USER	UFMPYLjtZJFuqcoeSBcimfTSLZqkUNoUVyLxxTTUiWmfWGYDeq
54	2020-04-11 09:50:54.248526+00	2020-04-11 09:50:54.248526+00	2020-04-11 10:01:43.739232+00	3	11	MYSQL_PASSWORD	6cc23f9b116db25bf383ee85a2bf4c84
55	2020-04-11 09:50:54.249559+00	2020-04-11 09:50:54.249559+00	2020-04-11 10:01:43.739232+00	3	11	MYSQL_DATABASE	UFMPYLjtZJFuqcoeSBci
56	2020-04-11 09:50:54.250557+00	2020-04-11 09:50:54.250557+00	2020-04-11 10:01:43.739232+00	3	11	MYSQL_ROOT_PASSWORD	6cc23f9b116db25bf383ee85a2bf4c84
57	2020-04-11 09:50:54.451616+00	2020-04-11 09:50:54.451616+00	2020-04-11 10:01:43.739232+00	3	11	MYSQL_HOST	svc-mysql-snowy-mountain.namespace-storm:3306
58	2020-04-11 10:08:13.992722+00	2020-04-11 10:08:13.992722+00	\N	3	12	MYSQL_ROOT_PASSWORD	5fbf7b77e17733eb2e8221e4113f306c
59	2020-04-11 10:08:13.997673+00	2020-04-11 10:08:13.997673+00	\N	3	12	MYSQL_USER	tmIeaAIDdmkqHhxrymkEdHZWtZFWXi
60	2020-04-11 10:08:13.998776+00	2020-04-11 10:08:13.998776+00	\N	3	12	MYSQL_PASSWORD	5fbf7b77e17733eb2e8221e4113f306c
61	2020-04-11 10:08:13.999694+00	2020-04-11 10:08:13.999694+00	\N	3	12	MYSQL_DATABASE	XqhgTCPcQMkWnMRzAFuE
62	2020-04-11 10:08:14.128162+00	2020-04-11 10:08:14.128162+00	\N	3	12	MYSQL_HOST	svc-mysql-snowy-mountain.namespace-storm:3306
68	2020-04-17 13:06:45.860517+00	2020-04-17 13:06:45.860517+00	\N	5	14	POSTGRES_USER	aUepmOLcNzqZJMYSJZWSMAPDptBJrxyliMY
69	2020-04-17 13:06:45.86129+00	2020-04-17 13:06:45.86129+00	\N	5	14	POSTGRES_PASSWORD	f2aa87bf4a80699657631088c5fdcf59
70	2020-04-17 13:06:45.867963+00	2020-04-17 13:06:45.867963+00	\N	5	14	POSTGRES_DB	aUepmOLcNzqZJMYSJZWS
71	2020-04-17 13:06:45.868746+00	2020-04-17 13:06:45.868746+00	\N	5	14	PG_DATA	/var/lib/postgresl/data/pg-bold-dew
72	2020-04-17 13:06:45.985244+00	2020-04-17 13:06:45.985244+00	\N	5	14	PG_HOST	svc-pg-bold-dew.namespace-storm:5432
73	2020-04-17 13:07:16.29504+00	2020-04-17 13:07:16.29504+00	\N	5	0	KEY1	VALUE1
63	2020-04-17 13:06:25.016889+00	2020-04-17 13:06:25.016889+00	2020-04-17 13:08:00.714048+00	5	13	MYSQL_USER	qviZakYXyTiMEnyHExMZHQpdHAVqRZ
64	2020-04-17 13:06:25.027641+00	2020-04-17 13:06:25.027641+00	2020-04-17 13:08:00.714048+00	5	13	MYSQL_PASSWORD	269804efb75cd44dc174fafe466fc148
65	2020-04-17 13:06:25.030059+00	2020-04-17 13:06:25.030059+00	2020-04-17 13:08:00.714048+00	5	13	MYSQL_DATABASE	pzpYdjoNgjWVJqDVkwQF
66	2020-04-17 13:06:25.035302+00	2020-04-17 13:06:25.035302+00	2020-04-17 13:08:00.714048+00	5	13	MYSQL_ROOT_PASSWORD	269804efb75cd44dc174fafe466fc148
67	2020-04-17 13:06:25.170744+00	2020-04-17 13:06:25.170744+00	2020-04-17 13:08:00.714048+00	5	13	MYSQL_HOST	svc-mysql-bold-dew.namespace-storm:3306
112	2020-04-19 19:27:30.61241+00	2020-04-19 19:27:30.61241+00	2020-04-19 19:33:08.663458+00	6	22	MONGO_INITDB_ROOT_USERNAME	ihJXwdYOMbXCEcomrziCYLJByUQBLDRKehCqFptvbObBKVLRVW
113	2020-04-19 19:27:40.769685+00	2020-04-19 19:27:40.769685+00	2020-04-19 19:33:08.663458+00	6	22	MONGO_HOST	svc-mongo-wild-cloud.namespace-storm:27017
110	2020-04-19 15:50:52.400596+00	2020-04-19 15:50:52.400596+00	2020-04-19 19:33:24.653191+00	6	0	BAR	FOO
89	2020-04-19 15:39:37.564908+00	2020-04-19 15:39:37.564908+00	2020-04-19 19:33:24.655639+00	6	0	FOO	BAR
114	2020-04-19 19:34:09.509521+00	2020-04-19 19:34:09.509521+00	\N	6	23	MONGO_INITDB_ROOT_PASSWORD	7a4ffbf73a93589392681981845a9523
74	2020-04-19 09:11:51.836533+00	2020-04-19 09:11:51.836533+00	2020-04-19 15:30:28.197085+00	6	15	POSTGRES_USER	PbPKTVZwgmjiOIFIqbiKGWSTeAdtEvvGqii
75	2020-04-19 09:11:51.846494+00	2020-04-19 09:11:51.846494+00	2020-04-19 15:30:28.197085+00	6	15	POSTGRES_PASSWORD	21f0851e0d63b94de189bddc8c8fa531
76	2020-04-19 09:11:51.848398+00	2020-04-19 09:11:51.848398+00	2020-04-19 15:30:28.197085+00	6	15	POSTGRES_DB	PbPKTVZwgmjiOIFIqbiK
77	2020-04-19 09:11:51.849058+00	2020-04-19 09:11:51.849058+00	2020-04-19 15:30:28.197085+00	6	15	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
78	2020-04-19 09:11:51.967357+00	2020-04-19 09:11:51.967357+00	2020-04-19 15:30:28.197085+00	6	15	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
115	2020-04-19 19:34:09.519513+00	2020-04-19 19:34:09.519513+00	\N	6	23	MONGO_INITDB_ROOT_USERNAME	nOmztkAYctSyxdyRrWauLsDULirGloyhKAooWBEAJGhSRHiATv
116	2020-04-19 19:34:19.674015+00	2020-04-19 19:34:19.674015+00	\N	6	23	MONGO_HOST	svc-mongo-wild-cloud.namespace-storm:27017
79	2020-04-19 15:30:44.591726+00	2020-04-19 15:30:44.591726+00	2020-04-19 15:37:54.668328+00	6	16	POSTGRES_DB	MyLufjPnhZNOtjPCTrgD
80	2020-04-19 15:30:44.598831+00	2020-04-19 15:30:44.598831+00	2020-04-19 15:37:54.668328+00	6	16	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
81	2020-04-19 15:30:44.599935+00	2020-04-19 15:30:44.599935+00	2020-04-19 15:37:54.668328+00	6	16	POSTGRES_USER	MyLufjPnhZNOtjPCTrgDcsPzFpkhigEkbbd
82	2020-04-19 15:30:44.601068+00	2020-04-19 15:30:44.601068+00	2020-04-19 15:37:54.668328+00	6	16	POSTGRES_PASSWORD	e0cb6d68cdb464ac81bfd401a7f8156e
83	2020-04-19 15:30:44.698077+00	2020-04-19 15:30:44.698077+00	2020-04-19 15:37:54.668328+00	6	16	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
84	2020-04-19 15:38:15.362676+00	2020-04-19 15:38:15.362676+00	2020-04-19 15:40:00.333946+00	6	17	POSTGRES_DB	kxTNhfkFxhnQjtPDQece
85	2020-04-19 15:38:15.369982+00	2020-04-19 15:38:15.369982+00	2020-04-19 15:40:00.333946+00	6	17	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
86	2020-04-19 15:38:15.371304+00	2020-04-19 15:38:15.371304+00	2020-04-19 15:40:00.333946+00	6	17	POSTGRES_USER	kxTNhfkFxhnQjtPDQecemKxLfxMaeFWEoCe
87	2020-04-19 15:38:15.372406+00	2020-04-19 15:38:15.372406+00	2020-04-19 15:40:00.333946+00	6	17	POSTGRES_PASSWORD	a8c1dd60ded0aa214c5564158085e43d
88	2020-04-19 15:38:18.526529+00	2020-04-19 15:38:18.526529+00	2020-04-19 15:40:00.333946+00	6	17	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
90	2020-04-19 15:40:19.999974+00	2020-04-19 15:40:19.999974+00	2020-04-19 15:44:21.000884+00	6	18	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
91	2020-04-19 15:40:20.001887+00	2020-04-19 15:40:20.001887+00	2020-04-19 15:44:21.000884+00	6	18	POSTGRES_USER	UIgODXiDlasMWshUrorxPhRnNHMRtFUyocw
92	2020-04-19 15:40:20.002956+00	2020-04-19 15:40:20.002956+00	2020-04-19 15:44:21.000884+00	6	18	POSTGRES_PASSWORD	3fbf4e2a78608e35351a687b7534d5ea
93	2020-04-19 15:40:20.004146+00	2020-04-19 15:40:20.004146+00	2020-04-19 15:44:21.000884+00	6	18	POSTGRES_DB	UIgODXiDlasMWshUrorx
94	2020-04-19 15:40:20.077071+00	2020-04-19 15:40:20.077071+00	2020-04-19 15:44:21.000884+00	6	18	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
95	2020-04-19 15:44:43.598626+00	2020-04-19 15:44:43.598626+00	2020-04-19 15:48:48.289243+00	6	19	POSTGRES_USER	lXadnJwzSzztpFxZhyMDBSuYiyzsoxcFgpZ
96	2020-04-19 15:44:43.603556+00	2020-04-19 15:44:43.603556+00	2020-04-19 15:48:48.289243+00	6	19	POSTGRES_PASSWORD	1820dbeed3cc6c5079f7b863128f29da
97	2020-04-19 15:44:43.604841+00	2020-04-19 15:44:43.604841+00	2020-04-19 15:48:48.289243+00	6	19	POSTGRES_DB	lXadnJwzSzztpFxZhyMD
98	2020-04-19 15:44:43.605804+00	2020-04-19 15:44:43.605804+00	2020-04-19 15:48:48.289243+00	6	19	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
99	2020-04-19 15:44:46.779471+00	2020-04-19 15:44:46.779471+00	2020-04-19 15:48:48.289243+00	6	19	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
100	2020-04-19 15:48:57.263855+00	2020-04-19 15:48:57.263855+00	2020-04-19 15:49:45.932522+00	6	20	POSTGRES_USER	kANMjXEjSfxBelvudYckIpFeGlNaFItYYSo
101	2020-04-19 15:48:57.268349+00	2020-04-19 15:48:57.268349+00	2020-04-19 15:49:45.932522+00	6	20	POSTGRES_PASSWORD	90439ee4a99a29abcd05d30c115425a9
102	2020-04-19 15:48:57.269434+00	2020-04-19 15:48:57.269434+00	2020-04-19 15:49:45.932522+00	6	20	POSTGRES_DB	kANMjXEjSfxBelvudYck
103	2020-04-19 15:48:57.270455+00	2020-04-19 15:48:57.270455+00	2020-04-19 15:49:45.932522+00	6	20	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
104	2020-04-19 15:49:01.433277+00	2020-04-19 15:49:01.433277+00	2020-04-19 15:49:45.932522+00	6	20	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
105	2020-04-19 15:50:01.918196+00	2020-04-19 15:50:01.918196+00	\N	6	21	POSTGRES_USER	sAMZqOTvzVqxOeqhBcrvVJeIvmhkVaTzQJG
106	2020-04-19 15:50:01.919657+00	2020-04-19 15:50:01.919657+00	\N	6	21	POSTGRES_PASSWORD	fdb98574aea682b51f39539d182e2f60
107	2020-04-19 15:50:01.920793+00	2020-04-19 15:50:01.920793+00	\N	6	21	POSTGRES_DB	sAMZqOTvzVqxOeqhBcrv
108	2020-04-19 15:50:01.921918+00	2020-04-19 15:50:01.921918+00	\N	6	21	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
109	2020-04-19 15:50:12.021971+00	2020-04-19 15:50:12.021971+00	\N	6	21	PG_HOST	svc-pg-wild-cloud.namespace-storm:5432
\.


--
-- Data for Name: plans; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.plans (id, created_at, updated_at, deleted_at, app_id, name, alias, price) FROM stdin;
3	2020-03-20 09:50:32.51456+00	2020-03-20 09:50:32.51456+00	\N	3	TEST	TST	0
4	2020-04-17 12:56:36.649434+00	2020-04-17 12:56:36.649434+00	\N	4	TEST	TST	0
5	2020-04-17 13:04:08.032712+00	2020-04-17 13:04:08.032712+00	\N	5	TEST	TST	0
6	2020-04-19 09:02:17.449899+00	2020-04-19 09:02:17.449899+00	\N	6	TEST	TST	0
7	2020-04-19 20:29:51.248107+00	2020-04-19 20:29:51.248107+00	\N	7	TEST	TST	0
8	2020-04-19 20:29:53.212122+00	2020-04-19 20:29:53.212122+00	\N	8	TEST	TST	0
9	2020-04-19 20:29:54.186756+00	2020-04-19 20:29:54.186756+00	\N	9	TEST	TST	0
10	2020-04-19 20:29:55.016188+00	2020-04-19 20:29:55.016188+00	\N	10	TEST	TST	0
11	2020-04-19 20:29:55.810082+00	2020-04-19 20:29:55.810082+00	\N	11	TEST	TST	0
\.


--
-- Data for Name: quota; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.quota (id, created_at, updated_at, deleted_at, memory, cpu, storage_size) FROM stdin;
\.


--
-- Data for Name: releases; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.releases (id, created_at, updated_at, deleted_at, app_id, last_check_sum, version_number, docker_url) FROM stdin;
1	2020-03-20 09:52:24.050592+00	2020-03-20 09:52:24.050592+00	\N	3		15	registry.csail.app/snowy-mountain:tsnXOjtA
2	2020-04-17 13:05:08.14687+00	2020-04-17 13:05:08.14687+00	\N	5		17	registry.csail.app/bold-dew:32986a
4	2020-04-19 11:27:15.276176+00	2020-04-19 11:27:15.276176+00	\N	6	\N	17	registry.csail.app/wild-cloud:54935e
5	2020-04-19 11:27:33.099104+00	2020-04-19 11:27:33.099104+00	\N	6	\N	18	registry.csail.app/wild-cloud:c85948
6	2020-04-19 11:27:48.250972+00	2020-04-19 11:27:48.250972+00	\N	6	\N	19	registry.csail.app/wild-cloud:04119b
3	2020-04-19 09:02:29.945256+00	2020-04-19 09:02:29.945256+00	\N	6		19	registry.csail.app/wild-cloud:04119b
7	2020-04-19 11:46:25.971699+00	2020-04-19 11:46:25.971699+00	\N	6	\N	20	registry.csail.app/wild-cloud:a3231c
8	2020-04-19 11:52:58.628115+00	2020-04-19 11:52:58.628115+00	\N	6	\N	20	registry.csail.app/wild-cloud:f3a670
9	2020-04-19 11:53:02.905003+00	2020-04-19 11:53:02.905003+00	\N	6	\N	20	registry.csail.app/wild-cloud:b806a0
10	2020-04-19 11:53:27.027804+00	2020-04-19 11:53:27.027804+00	\N	6	\N	20	registry.csail.app/wild-cloud:c5b930
11	2020-04-19 12:10:27.69701+00	2020-04-19 12:10:27.69701+00	\N	6	\N	20	registry.csail.app/wild-cloud:f1538d
12	2020-04-19 12:13:10.906495+00	2020-04-19 12:13:10.906495+00	\N	6	\N	20	registry.csail.app/wild-cloud:5aafd4
13	2020-04-19 19:40:56.408866+00	2020-04-19 19:40:56.408866+00	\N	6	\N	20	registry.csail.app/wild-cloud:311ce0
14	2020-04-19 19:44:26.936835+00	2020-04-19 19:44:26.936835+00	\N	6	\N	20	registry.csail.app/wild-cloud:d407ba
15	2020-04-19 19:44:30.564878+00	2020-04-19 19:44:30.564878+00	\N	6	\N	20	registry.csail.app/wild-cloud:434bfb
16	2020-04-19 19:44:35.31071+00	2020-04-19 19:44:35.31071+00	\N	6	\N	20	registry.csail.app/wild-cloud:f33874
17	2020-04-19 20:34:16.135843+00	2020-04-19 20:34:16.135843+00	\N	11	\N	1	registry.csail.app/muddy-shadow-129fcb2c:1c071c
18	2020-04-19 20:34:26.459366+00	2020-04-19 20:34:26.459366+00	\N	11	\N	2	registry.csail.app/muddy-shadow-129fcb2c:7dfc05
\.


--
-- Data for Name: resource_configs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resource_configs (id, created_at, updated_at, deleted_at, resource_id) FROM stdin;
1	2020-03-20 10:05:32.530541+00	2020-03-20 10:05:32.530541+00	\N	1
2	2020-03-22 09:18:16.145229+00	2020-03-22 09:18:16.145229+00	\N	2
3	2020-03-22 16:23:03.242418+00	2020-03-22 16:23:03.242418+00	\N	3
4	2020-03-22 16:27:53.568381+00	2020-03-22 16:27:53.568381+00	\N	4
5	2020-03-22 16:39:25.442821+00	2020-03-22 16:39:25.442821+00	\N	5
6	2020-03-22 16:56:27.919899+00	2020-03-22 16:56:27.919899+00	\N	6
7	2020-03-22 17:03:23.952229+00	2020-03-22 17:03:23.952229+00	\N	7
8	2020-04-10 19:43:41.322003+00	2020-04-10 19:43:41.322003+00	\N	8
9	2020-04-10 19:52:04.167699+00	2020-04-10 19:52:04.167699+00	\N	9
10	2020-04-10 19:57:18.121415+00	2020-04-10 19:57:18.121415+00	\N	10
11	2020-04-11 09:50:54.237874+00	2020-04-11 09:50:54.237874+00	\N	11
12	2020-04-11 10:08:13.990613+00	2020-04-11 10:08:13.990613+00	\N	12
13	2020-04-17 13:06:25.011112+00	2020-04-17 13:06:25.011112+00	\N	13
14	2020-04-17 13:06:45.860132+00	2020-04-17 13:06:45.860132+00	\N	14
15	2020-04-19 09:11:51.829511+00	2020-04-19 09:11:51.829511+00	\N	15
16	2020-04-19 15:30:44.588469+00	2020-04-19 15:30:44.588469+00	\N	16
17	2020-04-19 15:38:15.359866+00	2020-04-19 15:38:15.359866+00	\N	17
18	2020-04-19 15:40:19.999465+00	2020-04-19 15:40:19.999465+00	\N	18
19	2020-04-19 15:44:43.595498+00	2020-04-19 15:44:43.595498+00	\N	19
20	2020-04-19 15:48:57.260789+00	2020-04-19 15:48:57.260789+00	\N	20
21	2020-04-19 15:50:01.917593+00	2020-04-19 15:50:01.917593+00	\N	21
22	2020-04-19 19:27:30.606164+00	2020-04-19 19:27:30.606164+00	\N	22
23	2020-04-19 19:34:09.50628+00	2020-04-19 19:34:09.50628+00	\N	23
\.


--
-- Data for Name: resource_envs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resource_envs (id, created_at, updated_at, deleted_at, resource_id, env_key, env_value) FROM stdin;
1	2020-04-10 19:57:18.124553+00	2020-04-10 19:57:18.124553+00	\N	10	POSTGRES_USER	RnSSczXJdhcBBcNcRAmzmgjIozyOYdphFgf
2	2020-04-10 19:57:18.133168+00	2020-04-10 19:57:18.133168+00	\N	10	POSTGRES_PASSWORD	d89732d4198ca1da454d7d2aaf07b562
3	2020-04-10 19:57:18.134006+00	2020-04-10 19:57:18.134006+00	\N	10	POSTGRES_DB	RnSSczXJdhcBBcNcRAmz
4	2020-04-10 19:57:18.134754+00	2020-04-10 19:57:18.134754+00	\N	10	PG_DATA	/var/lib/postgresl/data/pg-snowy-mountain
5	2020-04-11 09:50:54.246803+00	2020-04-11 09:50:54.246803+00	\N	11	MYSQL_USER	UFMPYLjtZJFuqcoeSBcimfTSLZqkUNoUVyLxxTTUiWmfWGYDeq
6	2020-04-11 09:50:54.249056+00	2020-04-11 09:50:54.249056+00	\N	11	MYSQL_PASSWORD	6cc23f9b116db25bf383ee85a2bf4c84
7	2020-04-11 09:50:54.250065+00	2020-04-11 09:50:54.250065+00	\N	11	MYSQL_DATABASE	UFMPYLjtZJFuqcoeSBci
8	2020-04-11 09:50:54.251162+00	2020-04-11 09:50:54.251162+00	\N	11	MYSQL_ROOT_PASSWORD	6cc23f9b116db25bf383ee85a2bf4c84
9	2020-04-11 10:08:13.994534+00	2020-04-11 10:08:13.994534+00	\N	12	MYSQL_ROOT_PASSWORD	5fbf7b77e17733eb2e8221e4113f306c
10	2020-04-11 10:08:13.998222+00	2020-04-11 10:08:13.998222+00	\N	12	MYSQL_USER	tmIeaAIDdmkqHhxrymkEdHZWtZFWXi
11	2020-04-11 10:08:13.999212+00	2020-04-11 10:08:13.999212+00	\N	12	MYSQL_PASSWORD	5fbf7b77e17733eb2e8221e4113f306c
12	2020-04-11 10:08:14.000178+00	2020-04-11 10:08:14.000178+00	\N	12	MYSQL_DATABASE	XqhgTCPcQMkWnMRzAFuE
13	2020-04-17 13:06:25.02089+00	2020-04-17 13:06:25.02089+00	\N	13	MYSQL_USER	qviZakYXyTiMEnyHExMZHQpdHAVqRZ
14	2020-04-17 13:06:25.02877+00	2020-04-17 13:06:25.02877+00	\N	13	MYSQL_PASSWORD	269804efb75cd44dc174fafe466fc148
15	2020-04-17 13:06:25.031315+00	2020-04-17 13:06:25.031315+00	\N	13	MYSQL_DATABASE	pzpYdjoNgjWVJqDVkwQF
16	2020-04-17 13:06:25.035913+00	2020-04-17 13:06:25.035913+00	\N	13	MYSQL_ROOT_PASSWORD	269804efb75cd44dc174fafe466fc148
17	2020-04-17 13:06:45.860905+00	2020-04-17 13:06:45.860905+00	\N	14	POSTGRES_USER	aUepmOLcNzqZJMYSJZWSMAPDptBJrxyliMY
18	2020-04-17 13:06:45.867436+00	2020-04-17 13:06:45.867436+00	\N	14	POSTGRES_PASSWORD	f2aa87bf4a80699657631088c5fdcf59
19	2020-04-17 13:06:45.868377+00	2020-04-17 13:06:45.868377+00	\N	14	POSTGRES_DB	aUepmOLcNzqZJMYSJZWS
20	2020-04-17 13:06:45.869114+00	2020-04-17 13:06:45.869114+00	\N	14	PG_DATA	/var/lib/postgresl/data/pg-bold-dew
21	2020-04-19 09:11:51.841091+00	2020-04-19 09:11:51.841091+00	\N	15	POSTGRES_USER	PbPKTVZwgmjiOIFIqbiKGWSTeAdtEvvGqii
22	2020-04-19 09:11:51.847857+00	2020-04-19 09:11:51.847857+00	\N	15	POSTGRES_PASSWORD	21f0851e0d63b94de189bddc8c8fa531
23	2020-04-19 09:11:51.848775+00	2020-04-19 09:11:51.848775+00	\N	15	POSTGRES_DB	PbPKTVZwgmjiOIFIqbiK
24	2020-04-19 09:11:51.849468+00	2020-04-19 09:11:51.849468+00	\N	15	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
25	2020-04-19 15:30:44.595091+00	2020-04-19 15:30:44.595091+00	\N	16	POSTGRES_DB	MyLufjPnhZNOtjPCTrgD
26	2020-04-19 15:30:44.599397+00	2020-04-19 15:30:44.599397+00	\N	16	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
27	2020-04-19 15:30:44.600525+00	2020-04-19 15:30:44.600525+00	\N	16	POSTGRES_USER	MyLufjPnhZNOtjPCTrgDcsPzFpkhigEkbbd
28	2020-04-19 15:30:44.601636+00	2020-04-19 15:30:44.601636+00	\N	16	POSTGRES_PASSWORD	e0cb6d68cdb464ac81bfd401a7f8156e
29	2020-04-19 15:38:15.366882+00	2020-04-19 15:38:15.366882+00	\N	17	POSTGRES_DB	kxTNhfkFxhnQjtPDQece
30	2020-04-19 15:38:15.370678+00	2020-04-19 15:38:15.370678+00	\N	17	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
31	2020-04-19 15:38:15.371858+00	2020-04-19 15:38:15.371858+00	\N	17	POSTGRES_USER	kxTNhfkFxhnQjtPDQecemKxLfxMaeFWEoCe
32	2020-04-19 15:38:15.372954+00	2020-04-19 15:38:15.372954+00	\N	17	POSTGRES_PASSWORD	a8c1dd60ded0aa214c5564158085e43d
33	2020-04-19 15:40:20.000534+00	2020-04-19 15:40:20.000534+00	\N	18	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
34	2020-04-19 15:40:20.002411+00	2020-04-19 15:40:20.002411+00	\N	18	POSTGRES_USER	UIgODXiDlasMWshUrorxPhRnNHMRtFUyocw
35	2020-04-19 15:40:20.00363+00	2020-04-19 15:40:20.00363+00	\N	18	POSTGRES_PASSWORD	3fbf4e2a78608e35351a687b7534d5ea
36	2020-04-19 15:40:20.004687+00	2020-04-19 15:40:20.004687+00	\N	18	POSTGRES_DB	UIgODXiDlasMWshUrorx
37	2020-04-19 15:44:43.60069+00	2020-04-19 15:44:43.60069+00	\N	19	POSTGRES_USER	lXadnJwzSzztpFxZhyMDBSuYiyzsoxcFgpZ
38	2020-04-19 15:44:43.604251+00	2020-04-19 15:44:43.604251+00	\N	19	POSTGRES_PASSWORD	1820dbeed3cc6c5079f7b863128f29da
39	2020-04-19 15:44:43.60532+00	2020-04-19 15:44:43.60532+00	\N	19	POSTGRES_DB	lXadnJwzSzztpFxZhyMD
40	2020-04-19 15:44:43.606282+00	2020-04-19 15:44:43.606282+00	\N	19	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
41	2020-04-19 15:48:57.265874+00	2020-04-19 15:48:57.265874+00	\N	20	POSTGRES_USER	kANMjXEjSfxBelvudYckIpFeGlNaFItYYSo
42	2020-04-19 15:48:57.268911+00	2020-04-19 15:48:57.268911+00	\N	20	POSTGRES_PASSWORD	90439ee4a99a29abcd05d30c115425a9
43	2020-04-19 15:48:57.269941+00	2020-04-19 15:48:57.269941+00	\N	20	POSTGRES_DB	kANMjXEjSfxBelvudYck
44	2020-04-19 15:48:57.271131+00	2020-04-19 15:48:57.271131+00	\N	20	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
45	2020-04-19 15:50:01.91884+00	2020-04-19 15:50:01.91884+00	\N	21	POSTGRES_USER	sAMZqOTvzVqxOeqhBcrvVJeIvmhkVaTzQJG
46	2020-04-19 15:50:01.920232+00	2020-04-19 15:50:01.920232+00	\N	21	POSTGRES_PASSWORD	fdb98574aea682b51f39539d182e2f60
47	2020-04-19 15:50:01.921386+00	2020-04-19 15:50:01.921386+00	\N	21	POSTGRES_DB	sAMZqOTvzVqxOeqhBcrv
48	2020-04-19 15:50:01.922511+00	2020-04-19 15:50:01.922511+00	\N	21	PG_DATA	/var/lib/postgresl/data/pg-wild-cloud
49	2020-04-19 19:27:30.611201+00	2020-04-19 19:27:30.611201+00	\N	22	MONGO_INITDB_ROOT_PASSWORD	f0683cb808925faa147eb73d3a0f51b9
50	2020-04-19 19:27:30.613226+00	2020-04-19 19:27:30.613226+00	\N	22	MONGO_INITDB_ROOT_USERNAME	ihJXwdYOMbXCEcomrziCYLJByUQBLDRKehCqFptvbObBKVLRVW
51	2020-04-19 19:34:09.516395+00	2020-04-19 19:34:09.516395+00	\N	23	MONGO_INITDB_ROOT_PASSWORD	7a4ffbf73a93589392681981845a9523
52	2020-04-19 19:34:09.520142+00	2020-04-19 19:34:09.520142+00	\N	23	MONGO_INITDB_ROOT_USERNAME	nOmztkAYctSyxdyRrWauLsDULirGloyhKAooWBEAJGhSRHiATv
\.


--
-- Data for Name: resources; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.resources (id, created_at, updated_at, deleted_at, app_id, name) FROM stdin;
1	2020-03-20 10:05:32.526964+00	2020-03-20 10:05:32.526964+00	2020-03-22 09:18:05.510534+00	3	pg
2	2020-03-22 09:18:16.142109+00	2020-03-22 09:18:16.142109+00	2020-03-22 16:22:49.465669+00	3	pg
3	2020-03-22 16:23:03.239259+00	2020-03-22 16:23:03.239259+00	2020-03-22 16:27:40.8139+00	3	pg
4	2020-03-22 16:27:53.567718+00	2020-03-22 16:27:53.567718+00	2020-03-22 16:38:59.289577+00	3	pg
5	2020-03-22 16:39:25.439043+00	2020-03-22 16:39:25.439043+00	2020-03-22 16:56:05.390791+00	3	pg
6	2020-03-22 16:56:27.908122+00	2020-03-22 16:56:27.908122+00	2020-03-22 17:03:12.605369+00	3	pg
7	2020-03-22 17:03:23.951515+00	2020-03-22 17:03:23.951515+00	2020-04-10 19:43:22.288473+00	3	pg
8	2020-04-10 19:43:41.31523+00	2020-04-10 19:43:41.31523+00	2020-04-10 19:51:55.351413+00	3	pg
9	2020-04-10 19:52:04.167045+00	2020-04-10 19:52:04.167045+00	2020-04-10 19:57:13.568583+00	3	pg
10	2020-04-10 19:57:18.119249+00	2020-04-10 19:57:18.119249+00	\N	3	pg
11	2020-04-11 09:50:54.227016+00	2020-04-11 09:50:54.227016+00	2020-04-11 10:01:43.838046+00	3	mysql
12	2020-04-11 10:08:13.988888+00	2020-04-11 10:08:13.988888+00	\N	3	mysql
14	2020-04-17 13:06:45.859677+00	2020-04-17 13:06:45.859677+00	\N	5	pg
13	2020-04-17 13:06:25.004876+00	2020-04-17 13:06:25.004876+00	2020-04-17 13:08:00.75109+00	5	mysql
15	2020-04-19 09:11:51.826836+00	2020-04-19 09:11:51.826836+00	2020-04-19 15:30:28.243535+00	6	pg
16	2020-04-19 15:30:44.58636+00	2020-04-19 15:30:44.58636+00	2020-04-19 15:37:54.726407+00	6	pg
17	2020-04-19 15:38:15.357286+00	2020-04-19 15:38:15.357286+00	2020-04-19 15:40:00.369349+00	6	pg
18	2020-04-19 15:40:19.998791+00	2020-04-19 15:40:19.998791+00	2020-04-19 15:44:21.038563+00	6	pg
19	2020-04-19 15:44:43.589337+00	2020-04-19 15:44:43.589337+00	2020-04-19 15:48:48.359421+00	6	pg
20	2020-04-19 15:48:57.258516+00	2020-04-19 15:48:57.258516+00	2020-04-19 15:49:45.967597+00	6	pg
21	2020-04-19 15:50:01.916932+00	2020-04-19 15:50:01.916932+00	\N	6	pg
22	2020-04-19 19:27:30.603401+00	2020-04-19 19:27:30.603401+00	2020-04-19 19:33:08.760111+00	6	mongo
23	2020-04-19 19:34:09.501837+00	2020-04-19 19:34:09.501837+00	\N	6	mongo
\.


--
-- Name: accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.accounts_id_seq', 1, true);


--
-- Name: apps_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.apps_id_seq', 11, true);


--
-- Name: deployment_settings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.deployment_settings_id_seq', 4, true);


--
-- Name: domains_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.domains_id_seq', 1, true);


--
-- Name: environments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.environments_id_seq', 116, true);


--
-- Name: plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.plans_id_seq', 11, true);


--
-- Name: quota_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.quota_id_seq', 1, false);


--
-- Name: releases_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.releases_id_seq', 18, true);


--
-- Name: resource_configs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resource_configs_id_seq', 23, true);


--
-- Name: resource_envs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resource_envs_id_seq', 52, true);


--
-- Name: resources_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.resources_id_seq', 23, true);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: apps apps_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.apps
    ADD CONSTRAINT apps_pkey PRIMARY KEY (id);


--
-- Name: deployment_settings deployment_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deployment_settings
    ADD CONSTRAINT deployment_settings_pkey PRIMARY KEY (id);


--
-- Name: domains domains_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.domains
    ADD CONSTRAINT domains_pkey PRIMARY KEY (id);


--
-- Name: environments environments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.environments
    ADD CONSTRAINT environments_pkey PRIMARY KEY (id);


--
-- Name: plans plans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plans
    ADD CONSTRAINT plans_pkey PRIMARY KEY (id);


--
-- Name: quota quota_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.quota
    ADD CONSTRAINT quota_pkey PRIMARY KEY (id);


--
-- Name: releases releases_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.releases
    ADD CONSTRAINT releases_pkey PRIMARY KEY (id);


--
-- Name: resource_configs resource_configs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_configs
    ADD CONSTRAINT resource_configs_pkey PRIMARY KEY (id);


--
-- Name: resource_envs resource_envs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resource_envs
    ADD CONSTRAINT resource_envs_pkey PRIMARY KEY (id);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: idx_accounts_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_accounts_deleted_at ON public.accounts USING btree (deleted_at);


--
-- Name: idx_apps_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_apps_deleted_at ON public.apps USING btree (deleted_at);


--
-- Name: idx_deployment_settings_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_deployment_settings_deleted_at ON public.deployment_settings USING btree (deleted_at);


--
-- Name: idx_domains_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_domains_deleted_at ON public.domains USING btree (deleted_at);


--
-- Name: idx_environments_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_environments_deleted_at ON public.environments USING btree (deleted_at);


--
-- Name: idx_plans_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_plans_deleted_at ON public.plans USING btree (deleted_at);


--
-- Name: idx_quota_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_quota_deleted_at ON public.quota USING btree (deleted_at);


--
-- Name: idx_releases_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_releases_deleted_at ON public.releases USING btree (deleted_at);


--
-- Name: idx_resource_configs_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_resource_configs_deleted_at ON public.resource_configs USING btree (deleted_at);


--
-- Name: idx_resource_envs_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_resource_envs_deleted_at ON public.resource_envs USING btree (deleted_at);


--
-- Name: idx_resources_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_resources_deleted_at ON public.resources USING btree (deleted_at);


--
-- PostgreSQL database dump complete
--

