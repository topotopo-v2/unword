import type { Word } from '../types/word'
import SaveButton from "./SaveButton.tsx";
import { getCountryFlag } from '../utils/countryFlag'

type WordCardProps = {
    word: Word
    saved: boolean
    onSaveToggle: () => void
}

function WordCard({ word, saved, onSaveToggle }: WordCardProps) {
    return (
        <article className="word-card">
            <p className="word-date">{formatWordDate(word.word_date)}</p>

            <span className="country-flag">
                {getCountryFlag(word.country_code)}
            </span>

            <p className="language">
                {word.language}
            </p>

            {word.native_script && (
                <p className="native-script">{word.native_script}</p>
            )}

            {!word.native_script && (<h2 className="native-script">{word.word}</h2>)}

            <p className="pronunciation">{word.pronunciation}</p>

            <p className="definition">{word.definition}</p>

            <SaveButton
                saved={saved}
                onToggle={onSaveToggle}
            />

            {word.source && (
                <a
                    className="source-link"
                    href={word.source}
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Source
                </a>
            )}
        </article>
    )
}

// todo: MXN_ mobe to util
export function formatWordDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        timeZone: 'UTC',
    }).toUpperCase()
}

export default WordCard