--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8 (Debian 16.8-1.pgdg120+1)
-- Dumped by pg_dump version 16.8 (Debian 16.8-1.pgdg120+1)

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: blockchain_enum; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.blockchain_enum AS ENUM (
    'ethereum',
    'bitcoin',
    'solana'
);


ALTER TYPE public.blockchain_enum OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: blockchain_keys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.blockchain_keys (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id integer NOT NULL,
    name text,
    description text,
    blockchain public.blockchain_enum NOT NULL,
    network text NOT NULL,
    address text NOT NULL,
    encrypted_key text NOT NULL,
    public_key text NOT NULL,
    mnemonic_hash text NOT NULL,
    salt text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_deleted boolean DEFAULT false,
    deleted_at timestamp without time zone
);


ALTER TABLE public.blockchain_keys OWNER TO postgres;

--
-- Name: blockchain_keys_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.blockchain_keys_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.blockchain_keys_user_id_seq OWNER TO postgres;

--
-- Name: blockchain_keys_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.blockchain_keys_user_id_seq OWNED BY public.blockchain_keys.user_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    master_password text,
    login text NOT NULL,
    is_verified boolean DEFAULT false,
    is_active_master boolean DEFAULT false,
    is_active boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_deleted boolean DEFAULT false,
    deleted_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: blockchain_keys user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blockchain_keys ALTER COLUMN user_id SET DEFAULT nextval('public.blockchain_keys_user_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: blockchain_keys; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.blockchain_keys (id, user_id, name, description, blockchain, network, address, encrypted_key, public_key, mnemonic_hash, salt, created_at, updated_at, is_deleted, deleted_at) FROM stdin;
71613cac-1e70-4797-9cc8-d94d1e49bfd2	1	Test1	description Test1	ethereum	goerli	0x227A3620a76d35dC90B4b63eE39c5aD962B3E994	nVyyoFYaaOW3Sq6u1HixJUTLOtN/5x0dMV+z5x7MPON7wmWIJuDyLqeaAumgGZphEnmT7ZenxNB/vrDu	04f4d40856d3ab2f7d7e2abbb5ad6c6760523d5bbb8197e84860d07492b74f3d38ed62393453572954b862c78252c0d706513cccfd19e14513a9de790f31fa70f2	fa9afa2a4945c09ea6d1de5c88155246d9d40cffcc0cda775b4b3c3e267a077d	SntpVC86hUiI54RPy6mwMjHzS5Z0NA13xmxueeT/tlc=	2025-03-22 16:36:28.974797	2025-03-22 16:36:28.974797	f	\N
6c4b8a97-8b36-4958-9a45-68b380378364	1	dfvdfvd	dfvdfvdfv	ethereum	goerli	0x102483A8e55CE4a4E29a5ADC7814ACA7ED1dEbeD	rsw+RIo2c+zns6aes5rXtNrmOdutJEJINiSA//k5vGLYj1e1qC2yHs2Clh1fDn8Vt19Sf3FC/Gl+Bl8J	0452381e8c2da72663f5cacae838d9e2792252816014d42f9955df6b02491a2ebaec131125af916df2170064c9dad2e259a821a3fdd00159887bfd984485f4491a	525a07003a3f9b1a9e33820cfb79223d17bdd015805cee8dae1e0d564e9d67fc	rynfDYDn655y5KR/kxkffoEtSDsKsppRe3vw5PdbpPQ=	2025-03-24 01:39:48.262363	2025-03-24 01:39:48.262363	f	\N
7d9c3a37-c0bd-4f42-955d-653e26d96a39	1	cvbcvbc	cvbcvb	ethereum	goerli	0x14598519aCEa942825144f470a99bc136fed1400	8Ru1/thsLRt8wm3JCaSCdynyDyNGjmIPjNfgUCdnNxObaxaETEQdPuIzPsSV7QCf9oW6X2rxmi0LWDdU	04b2a75b87d68786c6fc187ec6bcc348d7959c28a5874fab7fcf9b9505e2ff1511e4907260020d2e7ce002bd464c5d745ed0f5e44da8d3ecfdce43ba89b159b2ab	b2adaf65ca5cbf5deee351ea1ad6a1dbb50bd08e727c026030e920b35372ec11	adb4VmgRxSySPlme3dGosd+PTbpAwbZyLSL1DRCBQRQ=	2025-03-25 20:10:44.950849	2025-03-25 20:10:44.950849	f	\N
edd97522-0f94-4b93-90fd-490bc22278cd	1	123344345	fbfgbfgbf	ethereum	goerli	0x85A3Afbe0B9ff89Fc6Ab16055eb49835F6A7adb9	4dz4FqBcG20hYhx1XPKh3SKjSoj7GZyJIyE42l4kODshVOgUbpyrob1Tb3FszE3d3kTK4SJ7LuNoEprf	045547f1c46f1c8a840eecaca0f2df681785d2e02c3936a75f29b210b962bc439e3af29fb45199d85621f1e230e7ad746a0d56d5004362384e52323dff47148038	a18d44dce54fc47789194de523daa097747530ea21dc4cb72fa0f74d9d15f6b5	Pk4e1Zw9pp82MCX6G5eazmPUOL8M6SPUjauBPCI3rjY=	2025-03-28 21:49:19.152682	2025-03-28 21:49:19.152682	f	\N
5f8aa20a-4487-4f10-a2f9-3977c0a6595e	1	xfhdfhdrht	rdthdthdrth	ethereum	goerli	0xe5b8D40483Fab6F9D439cE4C0E11a74F28d9A805	KhXTYlVfswpRDR3do6DMYhEOp2GWDTSU0plWZuJXYD7jrvMihsA4wetPGst9VzXj964E71cl1iKQvvEq	04909153712176d7c5c0e42e14075e0d364f563dc9afeb89980065d6939382c9d8efb29cc0626bf46c6a172d7fadbd0c9644d12c6a0b77c480c0a1fb303d0c0c84	91487c9ee1fbfa6377b0c0a884427123aff7858c98e7119831e378ae1626ef0f	qIWuyA6lQlQaW9igk3wS8tb9D8ofdmgm/Bt5xugiROg=	2025-03-28 22:01:27.472782	2025-03-28 22:01:27.472782	f	\N
233a0007-2bc1-44ea-a90d-11ea84d10950	1	Test Import	Test	ethereum	goerli	0xa97A56FB44e30304a4aF40E6Ca91467F3CdFB1b5	K/JYQFc0MYbp5OR56JNFM/NHPTSi1Ucd5Kht7iAAb471tr+xXrlYR7ScYPHDH/fosurKBCOC4o9MQajs	045abd6ea2e8910f376c1b65b421cda6ec0ff221fb775c087b6fc07036efff303c6df7ff8e581dbc1f3de58fa99926631f085a5d4f5b1a7931801c72b8516fd2d6	5f54227b74fbba7743c47cd286b4873f2e17331518d56facfc03e34cde4a0950	hC7wyQ1ZQFYPDo85yO3SlmqQc68AkkRiXPI/mZUst2M=	2025-06-01 17:28:19.463711	2025-06-01 17:28:19.463711	f	\N
6e01bcba-85e5-43ed-866e-265ea2ed8a0f	1	Test Import 2	Test 2	ethereum	goerli	0x39D20Ae8281204Fc263e580c706F89e63EC5C0B3	D9fNewILZNEV/8hMTUVRPhwjIFkFMFmulHw59eEjOTYUwmUVP0VX8Ua+hp2C58jYjRE/tLmhiQURqIhp	04f4e2519db17499aaf01525191127745120568017a5ed91557d8fa330d55c5e16eb3c7d64cddb71903887d9059cc58a57811e36ec217c22b02d8f5f8384ed5b27	7aa9e14a71be3cf31d9e3fe84f4c5123d55634b8851b967c9c3076fcb9a7e9af	d8f7mOucwy2cWf3PFTVO0QMFthGheI7t1SVWI2Gcakw=	2025-06-01 17:50:51.892613	2025-06-01 17:50:51.892613	f	\N
52ee4dc9-d700-49c4-8268-8710fcc5ac4c	1	Test mnemonic	mnemonic	ethereum	goerli	0xA1284C7826e14d773053F5837e9627a0CF80b214	rgB52jia85cdUerCqHuFf8KNNTva0U82UW0/XZYM2CZH519De22ScaQvWO3XA84rTyDDeyieM7YyYNUh	042f9358f587ad295f59734a93cb1d70d5c3bfe5d79cc3eb8ff84ca9378e0e9eb6a08b83a4e5f2b1f9460aeac3c292b3ae6b3de05eb5251a9b2e6997a35040949c	4125ab2785d16d8363938eafa707f7a0f90e1cf3128f24a1439c3ad9e7dc3329	tw46KoBE3rA2YERNFsd4K/7VrnmSQGLQa+VrYLgW9GA=	2025-06-17 09:25:23.929523	2025-06-17 09:25:23.929523	f	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, master_password, login, is_verified, is_active_master, is_active, created_at, updated_at, is_deleted, deleted_at) FROM stdin;
3	samgar.robot@gmail.com	cpG3Fa7kiyU8aPgg7GtR/Q==$xgMmWlA7EVmgNqh8OxVmGwyya32MmgHPl3WDQBN23A4=	\N	test_login	f	f	f	2025-06-17 10:02:14.340377	2025-06-17 10:02:14.340377	f	\N
1	samgarqwerty@gmail.com	nUaJxmbye6/LoLmyLum3hg==$8c/yO9bMxUTUhG5HeAvD8+bBi7B1ZzU3aNFcsdTzCSk=	7dRJh22qT68IlSKyVIj0Pw==$jqia9+0+ZcG8ej+OeopGFEhHbIgVk+87+euzRqzrtEc=	Samgar	t	t	t	2025-03-22 15:52:06.092449	2025-03-22 15:52:06.092449	f	\N
\.


--
-- Name: blockchain_keys_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.blockchain_keys_user_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: blockchain_keys blockchain_keys_address_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blockchain_keys
    ADD CONSTRAINT blockchain_keys_address_key UNIQUE (address);


--
-- Name: blockchain_keys blockchain_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blockchain_keys
    ADD CONSTRAINT blockchain_keys_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_blockchain_keys_address; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_blockchain_keys_address ON public.blockchain_keys USING btree (address);


--
-- Name: idx_blockchain_keys_blockchain; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_blockchain_keys_blockchain ON public.blockchain_keys USING btree (blockchain);


--
-- Name: idx_blockchain_keys_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_blockchain_keys_user_id ON public.blockchain_keys USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: blockchain_keys blockchain_keys_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.blockchain_keys
    ADD CONSTRAINT blockchain_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

