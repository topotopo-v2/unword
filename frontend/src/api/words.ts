import type { Word } from '../types/word'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL

export async function getTodayWord(timezone: string): Promise<Word> {
    const response = await fetch(
        `${API_BASE_URL}/api/words/today?timezone=${encodeURIComponent(timezone)}`
    )

    if (!response.ok) {
        throw new Error(`Failed to fetch today's word: ${response.status}`)
    }

    return response.json()
}

export async function getSavedWords(ids: string[]): Promise<Word[]> {
    if (ids.length === 0) {
        return []
    }

    const response = await fetch(
        `${API_BASE_URL}/api/words?ids=${encodeURIComponent(ids.join(','))}`
    )

    if (!response.ok) {
        throw new Error(`Failed to fetch saved words: ${response.status}`)
    }

    return response.json()
}