const STORAGE_KEY = 'unword-saved-words'

export function getSavedWordIds(): string[] {
    const value = localStorage.getItem(STORAGE_KEY)

    if (!value) {
        return []
    }

    return JSON.parse(value)
}

export function saveWord(id: string): void {
    const ids = getSavedWordIds()

    if (!ids.includes(id)) {
        ids.push(id)
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(ids))
}

export function unsaveWord(id: string): void {
    const ids = getSavedWordIds()
    const updatedIds = ids.filter(savedId => savedId !== id)

    localStorage.setItem(STORAGE_KEY, JSON.stringify(updatedIds))
}