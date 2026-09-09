import type {Word} from "../types/word.ts";
import {getCountryFlag} from "../utils/countryFlag.ts";

type SavedWordCardProps = {
    word: Word
    onSaveToggle: (id: string) => void
}

export function SavedWordCard({word, onSaveToggle}: SavedWordCardProps) {
    return (
        <article className="saved-word-card">
            <button
                className="unsave-button"
                onClick={() => onSaveToggle(word.id)}
                aria-label={`Unsave ${word.word}`}
            >
                <img src="/icons/heart-black.svg" alt=""/>
            </button>

            <p className="numeric-day">{getDateDay(word.word_date)}</p>

            <span className="line"/>

            <span className="saved-country-flag">
                {getCountryFlag(word.country_code)}
            </span>

            <div className="saved-word-details">
                {word.native_script && (
                    <p className="saved-native-script">
                        {word.native_script}
                    </p>
                )}

                <p className="saved-word">{word.word}</p>

                <p className="saved-definition">
                    {word.definition}
                </p>

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
            </div>
        </article>
    )
}

export function getDateDay(date: string): string {
    return new Date(date).toLocaleDateString('en', {
        day: 'numeric',
        timeZone: 'UTC',
    })
}