type NavigationProps = {
    currentPage: 'today' | 'saved'
    onNavigate: (page: 'today' | 'saved') => void
}

function Navigation({ currentPage, onNavigate }: NavigationProps) {
    return (
        <nav className="navigation">
            <button
                className={currentPage === 'today' ? 'active' : ''}
                onClick={() => onNavigate('today')}
            >
                Today's
            </button>

            <button
                className={currentPage === 'saved' ? 'active' : ''}
                onClick={() => onNavigate('saved')}
            >
                Saved
            </button>
        </nav>
    )
}

export default Navigation