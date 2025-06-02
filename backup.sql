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
8b78b9a2-77cf-48e9-99c1-aaa53ab44fc5	1	Test	rgergerg	ethereum	goerli	0xe4fab3cBE0241186a8e2DDC802B3D6e317E2f1d7	YDRw+zTKrdUl5p97vARB2uyPGTK9rK/+/bFp0voallAoHq49Lp36OnLLze9S0Fs49KUGlPfBmElcaac2	046ff680ce25b6a9bd3eb04d87c583fd77b6c36005281adb95307536d61680dd87581d122576a8736d8a7ccce7d04b5ca7210f0b4ae22cecc2363edfb51df2ddef	6b9f25226e60a909b8d589999f020df3d3c45e14fb9f08e21c830060b02bef45	6ZhuETh9udmDnQKURd/WtocYG0WYVeY6+QQ0qEFON64=	2025-03-24 01:00:19.689678	2025-03-24 01:00:19.689678	f	\N
a6117eed-404d-4d5b-89cf-a0285d59725e	1	kjkjdf	dkjfbvdkjf	ethereum	goerli	0x77Dd99B9c503379f817a817bE03304c1047472F7	sgl5Ud2HV0Z6BqpVkK8pVr/jfmcjqFgI6sGqXVJyVIxHWCZm6mxxP8FWUQ7zG4pxx5dIAOqfAWbT5fe0	049576ece40e10015a7bdca159af5e583c8fe59958fc599593237bc926239cd49e2d8bd7a7361e25e4785f4852770451588b29fd9ea8e926668d7fb9dfad23e9e0	256b11042bd8a243ccb9f3d9e2cccecc276cdf4b56c713572f2a0850dd16be1d	qjl0IANiiITJ0P2/R9BMNzcNcLNJBospT+pe4HLm18I=	2025-03-24 01:25:54.83512	2025-03-24 01:25:54.83512	f	\N
ece79333-a308-4157-b726-bdd36a1033cd	1	кнерекнркр	керкеркер	ethereum	goerli	0x5f0a73A0Eb08267D70E41a086C1dd935A1e41dd5	J9QPW8/IrkRdhcsXxXwBr1adleIi7bjKdJucxU6RaS8F15zM0bs8oeRk50WlEn+dauhPLXDNHSGkEJv1	04dd64eea34cc1fbb37c9a335548080892ea82d9a0ee36980576ceae8c132db470f741eeb8a7b42d7f743d7c0525f4d02c96ba5273a82a29c04b5c937bd9b9e99b	7f53b0b862e933192647716188e3658e6639207d37a24855e458c5f5bddceeac	TYvMLrDvpHQAx6DLrEXWTokqjNGObtYSKHo/q8spNGY=	2025-03-24 01:28:25.762209	2025-03-24 01:28:25.762209	f	\N
bf839721-7aac-4486-8aa9-a4cfa4c27bce	1	rfergerg	ergergerg	ethereum	goerli	0xAC140DD1F4Cd3dB674A41e941Dee3Eda66DDABb5	3HNEKuI6tOoaXUXBk9xDlqmYpwxw9Yr8H7bWx1rrnmE9OjwZb1AAyqvw4mm31btQ/ng5eBVU060z52L4	044efe5f1fdee8d3e3d8e1441ade81d40f59b8b1b8d9730248d448f6c93568aec42c2fe5921ee327cd565844ffd2f340836a9ca5db0d2d490f37616eb9d053d939	0b85a324434aa8bbb6817f88aa87caf1e446c4778de922b81755f83a4c5c32fc	sbEzuyD04dnIHUSxEwaraVBlkgP/+R7iQFGn1Rt0lwU=	2025-03-24 01:29:06.653451	2025-03-24 01:29:06.653451	f	\N
912066ba-bf12-4430-8fd8-dc69e4cbea76	1	rgefdfbdf	dfvdfvdfv	ethereum	goerli	0x3d9E0b44889C931b5945210e17a0c02de2f3D856	kRrStuio5IXatZ1dPRAFGSaBpdAuPzLAV3fmmwVc+r16HicCWhUwiA58yvVnQjNZWQeTR0zFyqg35Fqe	040248c94fb71fe3f3b473f5028b3f641bda54821f0d00b1edcf877228e6cfb2bb9d95c616b5635eec29a7463fd54468ccf53069ff4447186e85c96ecd3077d99f	5306645b420cc5e1cabe712cd94c751907c3de6defdb30e54b1e0b196cc03e35	Vsqf5IbeEudbzKZgqBu5vT6oAorEO77GXMW60iKLMI4=	2025-03-24 01:29:49.663766	2025-03-24 01:29:49.663766	f	\N
bb3b65ff-edd1-4e3d-a516-08dda8495874	1	fvdfvdfvdfv	dfvdfvdfvdfv	ethereum	goerli	0x6a877BAd4456417CC8BF235285420f879C77ef86	2szY2Coo3OOctYDwotKxSuFaDObYWpmi2R7lfSX2gv4FemV4OYMErx1MYVeMk6XI1+JmYqSQNKhPbt6H	04f5d8e3759217e8ba0b3bef7e0b14857273e477b2a45284a1db621f54bf27ebf10482ebc32ecef46a02d04b5d4eb44a8623b3237b74a28f84761b6cd974f83f4a	5bd6a574c8453cbe0d6f37667049e268b0e6f0135219302e0225d31538de6a0c	/L4V1368tXEqnyvtVDJ8izrGI6ObvHsC8RO+INKATUc=	2025-03-24 01:33:59.32242	2025-03-24 01:33:59.32242	f	\N
19fce58a-e69f-40e3-9ec7-3f83bd606116	1	dfvdfvdfv	dfvdfvdfv	ethereum	goerli	0xdC8c22Fd036B77f98558082EF4430Ae9e9049780	gancXxFwfB3d01vLO7xhVlmOhhWrVO+ToQbHUlhqA3OdKQPoMgZAdaaEDrdSFP6DVFyRbK81uykeEbSu	040890a51a16c10999b494d31f19cbd35f7d549a93f79a7f6a4c42d9452c6ba696a1c7e85b1acfefae1e745292079631fdfb580798099314b04fbac29a0672f020	acc2636010153dee7f31b31a47c85d84ac5ffbbc14e41fb38832d8172dc95814	zg3hHF4Z3VfOeJHyg69qEx29wHrjirICcZEcHKvigIo=	2025-03-24 01:37:00.723268	2025-03-24 01:37:00.723268	f	\N
2add8691-b2b0-4953-b533-564379c76c15	1	dfvdfvd	dfvdfvdfv	ethereum	goerli	0x0A6c62C4259825687c8465c3aA650f8702094D36	If5kAdFIWswwB0YP1xk5+g97/tB/+8Q9bI0t/5bohU0RbfMzXG/eDmqRMi2DGkhBcsUpHntH5iZIOC2E	043fa73eb7ffd77ab119f1d2c3aef57e918f08179ea102b80a5d638eb134c11f60e3eb855a303ae848adcd9884db1d6c058bf3f9ca8e6a69c89501fdc5878be6e1	a7c335bd93b16ef936f5ac195664296340379c6165cade07b21d66db4ccb2ffe	LpdW8xZbSByWdU4o7XE7juNabsCrxvkFvkCZf9LoGNo=	2025-03-24 01:38:39.522425	2025-03-24 01:38:39.522425	f	\N
6c4b8a97-8b36-4958-9a45-68b380378364	1	dfvdfvd	dfvdfvdfv	ethereum	goerli	0x102483A8e55CE4a4E29a5ADC7814ACA7ED1dEbeD	rsw+RIo2c+zns6aes5rXtNrmOdutJEJINiSA//k5vGLYj1e1qC2yHs2Clh1fDn8Vt19Sf3FC/Gl+Bl8J	0452381e8c2da72663f5cacae838d9e2792252816014d42f9955df6b02491a2ebaec131125af916df2170064c9dad2e259a821a3fdd00159887bfd984485f4491a	525a07003a3f9b1a9e33820cfb79223d17bdd015805cee8dae1e0d564e9d67fc	rynfDYDn655y5KR/kxkffoEtSDsKsppRe3vw5PdbpPQ=	2025-03-24 01:39:48.262363	2025-03-24 01:39:48.262363	f	\N
5e00120a-8cb4-4f53-b3eb-20d14074757e	1	jgjyguy	jygujyhg	ethereum	goerli	0x9aAf5D989444EEeCc20DE924D2BFBAfEcBcE6A22	eKYsJSdltQOg7jO7u7NvwfGla5j0v1BuM5fmE+aaBpXE6B1B8TwE1geckGvYxl+siprpKt4ycRENW/ne	04ad85e50b6ab34d30e8807d410a44223a695ff05f2b2a374d281e0b91734928eeecd110dcfe2dbe053035668ee39005abe01b415710c932b5addab8290b9dd49b	0b73a31370694b184f57528013bc5c82b156b1e40a33d24b2d557cc1257be069	pxD46MYMKgo6SrI56Vx103Qn/kCU+AeHERbj+waCe48=	2025-03-24 01:41:15.546021	2025-03-24 01:41:15.546021	f	\N
7d9c3a37-c0bd-4f42-955d-653e26d96a39	1	cvbcvbc	cvbcvb	ethereum	goerli	0x14598519aCEa942825144f470a99bc136fed1400	8Ru1/thsLRt8wm3JCaSCdynyDyNGjmIPjNfgUCdnNxObaxaETEQdPuIzPsSV7QCf9oW6X2rxmi0LWDdU	04b2a75b87d68786c6fc187ec6bcc348d7959c28a5874fab7fcf9b9505e2ff1511e4907260020d2e7ce002bd464c5d745ed0f5e44da8d3ecfdce43ba89b159b2ab	b2adaf65ca5cbf5deee351ea1ad6a1dbb50bd08e727c026030e920b35372ec11	adb4VmgRxSySPlme3dGosd+PTbpAwbZyLSL1DRCBQRQ=	2025-03-25 20:10:44.950849	2025-03-25 20:10:44.950849	f	\N
f9de3a95-5b30-41a1-85d1-af8e57879ce5	1	test2323	srkfjaifuhaef	ethereum	goerli	0xa97A56FB44e30304a4aF40E6Ca91467F3CdFB1b5	Z17n+Flw+U5YTruZrqqIU4mEgNtbJbgdaK9tYLDakpYm9wfRf/o9ipt0dYurMlBK6S6uasg1e3uWrVrQ	045abd6ea2e8910f376c1b65b421cda6ec0ff221fb775c087b6fc07036efff303c6df7ff8e581dbc1f3de58fa99926631f085a5d4f5b1a7931801c72b8516fd2d6	8641453a068867d951dfe9579ebde21bdc3f9f5cb361a51388cf5314cb6ed4d3	vOdpR7aWNhqhZbDro/8Q4lpkQ+MIGpRELpk5oU2EmgM=	2025-03-28 21:48:33.382462	2025-03-28 21:48:33.382462	f	\N
edd97522-0f94-4b93-90fd-490bc22278cd	1	123344345	fbfgbfgbf	ethereum	goerli	0x85A3Afbe0B9ff89Fc6Ab16055eb49835F6A7adb9	4dz4FqBcG20hYhx1XPKh3SKjSoj7GZyJIyE42l4kODshVOgUbpyrob1Tb3FszE3d3kTK4SJ7LuNoEprf	045547f1c46f1c8a840eecaca0f2df681785d2e02c3936a75f29b210b962bc439e3af29fb45199d85621f1e230e7ad746a0d56d5004362384e52323dff47148038	a18d44dce54fc47789194de523daa097747530ea21dc4cb72fa0f74d9d15f6b5	Pk4e1Zw9pp82MCX6G5eazmPUOL8M6SPUjauBPCI3rjY=	2025-03-28 21:49:19.152682	2025-03-28 21:49:19.152682	f	\N
5f8aa20a-4487-4f10-a2f9-3977c0a6595e	1	xfhdfhdrht	rdthdthdrth	ethereum	goerli	0xe5b8D40483Fab6F9D439cE4C0E11a74F28d9A805	KhXTYlVfswpRDR3do6DMYhEOp2GWDTSU0plWZuJXYD7jrvMihsA4wetPGst9VzXj964E71cl1iKQvvEq	04909153712176d7c5c0e42e14075e0d364f563dc9afeb89980065d6939382c9d8efb29cc0626bf46c6a172d7fadbd0c9644d12c6a0b77c480c0a1fb303d0c0c84	91487c9ee1fbfa6377b0c0a884427123aff7858c98e7119831e378ae1626ef0f	qIWuyA6lQlQaW9igk3wS8tb9D8ofdmgm/Bt5xugiROg=	2025-03-28 22:01:27.472782	2025-03-28 22:01:27.472782	f	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, master_password, login, is_verified, is_active_master, is_active, created_at, updated_at, is_deleted, deleted_at) FROM stdin;
1	samgarqwerty@gmail.com	RfhuXta4RlRPQeOjKDymHA==$Bv0BqZftEO7EGVznUobZhFUx7RRzAHtfVUa70KfrygI=	7dRJh22qT68IlSKyVIj0Pw==$jqia9+0+ZcG8ej+OeopGFEhHbIgVk+87+euzRqzrtEc=	Samgar	t	t	t	2025-03-22 15:52:06.092449	2025-03-22 15:52:06.092449	f	\N
\.


--
-- Name: blockchain_keys_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.blockchain_keys_user_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


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

