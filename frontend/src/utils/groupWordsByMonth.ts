import type { Word } from '../types/word'

export type WordGroup = {
    month: string
    words: Word[]
}

export function groupWordsByMonth(words: Word[]): WordGroup[] {
    const groups: WordGroup[] = []

    for (const word of words) {
        const month = new Date(word.word_date).toLocaleDateString('en-US', {
            month: 'short',
            year: 'numeric',
            timeZone: 'UTC',
        })

        const existingGroup = groups.find(group => group.month === month)

        if (existingGroup) {
            existingGroup.words.push(word)
        } else {
            groups.push({
                month,
                words: [word],
            })
        }
    }

    return groups
}