package resources

import (
	"log"

	"github.com/jackc/pgx"
)

/**
Vote queries
 */

//
//INSERTING
//

const putVoteByThrID = `WITH sub AS (
	INSERT INTO vote (user_id, thread_id, voice)
	VALUES (
		(SELECT id FROM users WHERE nickname=$1),
		$2, $3)
	ON CONFLICT ON CONSTRAINT unique_user_and_thread
		DO UPDATE
			SET prev_voice = vote.voice ,
				voice = EXCLUDED.voice
	RETURNING prev_voice, voice, thread_id)
UPDATE thread
SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub)
WHERE id = $2
RETURNING id,
	slug::TEXT,
	title, message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count`

const putVoteByThrSLUG = `WITH sub AS (
	INSERT INTO vote (user_id, thread_id, voice)
	VALUES (
		(SELECT id FROM users WHERE nickname=$1),
		(SELECT id FROM thread WHERE slug=$2),
		$3)
	ON CONFLICT ON CONSTRAINT unique_user_and_thread
	DO UPDATE
		SET prev_voice = vote.voice ,
			voice = EXCLUDED.voice
	RETURNING prev_voice,
		voice,
		thread_id)
UPDATE thread
SET votes_count = votes_count - (SELECT prev_voice-voice FROM sub)
WHERE slug=$2
RETURNING id,
	slug::TEXT,
	title,
	message,
	forum_slug::TEXT,
	user_nick::TEXT,
	created,
	votes_count`


//
// UPDATING
//

//
// CLAIMING
//

func PrepareVotesQureies(tx *pgx.ConnPool){
	if _, err := tx.Prepare("putVoteByThrID", putVoteByThrID); err != nil {
		log.Fatalln(err)
	}
	if _, err := tx.Prepare("putVoteByThrSLUG", putVoteByThrSLUG); err != nil {
		log.Fatalln(err)
	}
}