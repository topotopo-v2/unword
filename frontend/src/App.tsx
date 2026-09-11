import type {Word} from "./types/word.ts";
import WordCard from "./components/WordCard.tsx";
import {useEffect, useState} from "react";
import {getTodayWord} from "./api/words.ts";
import LoadingState from "./components/LoadingState.tsx";
import ErrorState from "./components/ErrorState.tsx";
import {getSavedWordIds, saveWord, unsaveWord} from "./storage/savedWords.ts";
import Header from "./components/Header.tsx";
import SavedWords from './pages/SavedWords';
import Navigation from './components/Navigation'


function App() {
    const [word, setWord] = useState<Word | null>(null)
    const [saved, setSaved] = useState(false)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [page, setPage] = useState<'today' | 'saved'>('today')

    async function loadWord() {
        setError(null)

        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone

        try {
            const todayWord = await getTodayWord(timezone)
            setWord(todayWord)

            const savedIds = getSavedWordIds()
            setSaved(savedIds.includes(todayWord.id))
        } catch (err) {
            console.error('Failed to load today\'s word:', err)
            setError("Unable to load today's word.")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        if (page === 'today' && word) {
            const savedIds = getSavedWordIds()
            setSaved(savedIds.includes(word.id))
        }
    }, [page, word])

    useEffect(() => {
        loadWord()
    }, [])

    function handleSaveToggle() {
        if (!word) {
            return
        }

        if (saved) {
            unsaveWord(word.id)
            setSaved(false)
        } else {
            saveWord(word.id)
            setSaved(true)
        }
    }

    return (
        <main>
            <Navigation
                currentPage={page}
                onNavigate={setPage}
            />
            <Header />

            {page === 'today' && (
                <>
                    {loading && <LoadingState message={"Loading today's word..."} />}
                    {!loading && error && (
                        <ErrorState
                            message={error}
                            onRetry={loadWord}
                        />
                    )}

                    {!loading && word && (
                        <WordCard
                            word={word}
                            saved={saved}
                            onSaveToggle={handleSaveToggle}
                        />
                    )}
                </>
            )}

            {page === 'saved' && (
                <SavedWords />
            )}
        </main>
    )

}

export default App
