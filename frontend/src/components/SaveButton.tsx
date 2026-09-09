type SaveButtonProps = {
    saved: boolean
    onToggle: () => void
}

function SaveButton({ saved, onToggle }: SaveButtonProps) {
    return (
        <button className="save-button"
            onClick={onToggle}
            aria-label={saved ? 'Unsave word' : 'Save word'}
        >
            <img src={saved ? "/icons/heart-black.svg" : "/icons/heart-outline.svg"}
            alt=""/>
            <span>{saved ? 'Saved' : 'Save'}</span>
        </button>
    )
}

export default SaveButton