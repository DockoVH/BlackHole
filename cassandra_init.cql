CREATE KEYSPACE IF NOT EXISTS blackhole 
WITH REPLICATION = {
	'class': 'SimpleStrategy',
	'replication_factor': 1
};

USE blackhole;

CREATE TABLE IF NOT EXISTS igraci (
	username TEXT PRIMARY KEY,
	lozinka_hes TEXT,
	datum_rodjenja DATE
);

CREATE TABLE IF NOT EXISTS sobe (
	uuid UUID PRIMARY KEY,
	kod TEXT,
	pobednik TEXT,
	vreme TIMESTAMP
);

CREATE TABLE IF NOT EXISTS igrac_sobe (
	username TEXT,
	soba_uuid UUID,
	PRIMARY KEY (username, soba_uuid)
) WITH CLUSTERING ORDER BY (soba_uuid DESC);

CREATE TABLE IF NOT EXISTS potezi (
	soba_uuid UUID,
	vreme TIMESTAMP,
	username TEXT,
	indeks_polja INT,
	PRIMARY KEY (soba_uuid, vreme)
) WITH CLUSTERING ORDER BY (vreme ASC);

CREATE TABLE IF NOT EXISTS cet_poruke (
	soba_uuid UUID,
	vreme TIMESTAMP,
	username TEXT,
	sadrzaj TEXT,
	PRIMARY KEY (soba_uuid, vreme)
) WITH CLUSTERING ORDER BY (vreme ASC);

CREATE TABLE IF NOT EXISTS nepozeljne_reci (
	rec TEXT PRIMARY KEY
);





INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('acrotomophilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('anal') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('anilingus') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('anus') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('apeshit') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('arsehole') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('ass') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('asshole') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('assmunch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('autoerotic') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('babeland') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bangbros') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bareback') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('barenaked') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bastard') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bastardo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bastinado') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bbw') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bdsm') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('beaner') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('beaners') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bestiality') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bimbos') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('birdlock') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bitch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bitches') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('blowjob') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('blumpkin') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bollocks') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bondage') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('boner') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('boob') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('boobs') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bukkake') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bulldyke') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bullshit') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('bunghole') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('busty') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('butt') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('buttcheeks') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('butthole') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('camgirl') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('camslut') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('camwhore') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('carpetmuncher') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('circlejerk') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('clit') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('clitoris') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('clusterfuck') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cock') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cocks') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('coprolagnia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('coprophilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cornhole') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('coon') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('coons') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('creampie') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cum') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cumming') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cunnilingus') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('cunt') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('darkie') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('daterape') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('deepthroat') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dendrophilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dick') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dildo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dingleberry') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dingleberries') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('doggiestyle') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('doggystyle') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dolcett') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('domination') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dominatrix') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('dommes') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('ecchi') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('ejaculation') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('erotic') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('erotism') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('escort') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('eunuch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('faggot') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fecal') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('felch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fellatio') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('feltch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('femdom') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('figging') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fingerbang') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fingering') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fisting') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('footjob') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('frotting') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fuck') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fuckin') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fucking') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fucktards') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('fudgepacker') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('futanari') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('genitals') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('goatcx') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('goatse') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('gokkun') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('goodpoop') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('goregasm') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('grope') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('guro') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('handjob') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('hardcore') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('hentai') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('homoerotic') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('honkey') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('hooker') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('humping') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('incest') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('intercourse') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('jailbait') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('jigaboo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('jiggaboo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('jiggerboo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('jizz') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('juggs') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('kike') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('kinbaku') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('kinkster') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('kinky') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('knobbing') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('lolita') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('lovemaking') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('masturbate') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('milf') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('motherfucker') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('muffdiving') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nambla') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nawashi') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('negro') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('neonazi') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nigga') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nigger') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nimphomania') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nipple') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nipples') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nude') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nudity') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nympho') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('nymphomania') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('octopussy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('omorashi') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('orgasm') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('orgy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('paedophile') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('paki') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('panties') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('panty') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pedobear') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pedophile') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pegging') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('penis') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pissing') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pisspig') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('playboy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('ponyplay') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('poof') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('poon') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('poontang') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('punany') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('poopchute') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('porn') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('porno') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pornography') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pthc') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pubes') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('pussy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('queaf') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('queef') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('quim') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('raghead') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('rape') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('raping') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('rapist') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('rectum') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('rimjob') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('rimming') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sadism') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('santorum') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('scat') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('schlong') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('scissoring') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('semen') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sex') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sexo') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sexy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shemale') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shibari') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shit') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shitblimp') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shitty') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shota') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('shrimping') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('skeet') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('slanteye') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('slut') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('smut') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('snatch') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('snowballing') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sodomize') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sodomy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('spic') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('splooge') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('spooge') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('spunk') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('strapon') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('strappado') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('suck') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('sucks') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('swastika') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('swinger') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('threesome') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('throating') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tit') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tits') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('titties') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('titty') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('topless') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tosser') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('towelhead') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tranny') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tribadism') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tubgirl') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('tushy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('twat') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('twink') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('twinkie') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('undressing') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('upskirt') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('urophilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('vagina') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('vibrator') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('vorarephilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('voyeur') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('vulva') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('wank') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('wetback') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('yaoi') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('yiffy') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('zoophilia') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('zvezda') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('kanta') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('korpa') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('radijator') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('krokodil') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('aligator') IF NOT EXISTS;
INSERT INTO blackhole.nepozeljne_reci (rec) VALUES ('zebra') IF NOT EXISTS;
