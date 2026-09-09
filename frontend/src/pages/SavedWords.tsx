import { useEffect, useState } from 'react'
import { getSavedWords } from '../api/words'
import type { Word } from '../types/word'
import { groupWordsByMonth } from '../utils/groupWordsByMonth'
import { getSavedWordIds, unsaveWord } from '../storage/savedWords'
import { SavedWordCard } from '../components/SavedWordCard'
import ErrorState from "../components/ErrorState"
import EmptyState from "../components/EmptyState.tsx";
import LoadingState from "../components/LoadingState.tsx";

function SavedWords() {
    const [words, setWords] = useState<Word[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)


    async function loadSavedWords() {
        setLoading(true)
        setError(null)

        try {
            const ids = getSavedWordIds()
            const savedWords = await getSavedWords(ids)

            setWords(savedWords)
        } catch (err) {
            console.error('Failed to load saved words:', err)
            setError('Unable to load saved words.')
        } finally {
            setLoading(false)
        }
    }


    useEffect(() => {
        loadSavedWords()
    }, [])

    function handleUnsave(id: string) {
        unsaveWord(id)

        setWords(currentWords =>
            currentWords.filter(word => word.id !== id)
        )
    }

    if (loading) {
        return <LoadingState message={"Loading..."} />
    }

    if (error) {
        return (
            <ErrorState
                message={error}
                onRetry={loadSavedWords}
            />
        )
    }

    if (words.length === 0) {
        return <EmptyState message="No saved words yet." />
    }

    const groups = groupWordsByMonth(words)

    return (
        <>
            {/*<p className="saved-words-title">YOUR SAVED WORDS</p>*/}

            {groups.map((group) => (
                <section key={group.month}>
                    <p className="saved-words-month">{group.month}</p>

                    {group.words.map((word) => (
                        <SavedWordCard
                            key={word.id}
                            word={word}
                            onSaveToggle={handleUnsave}
                        />
                    ))}
                </section>
            ))}
        </>
    )
}

export default SavedWords