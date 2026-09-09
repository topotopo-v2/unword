type ErrorStateProps = {
    message: string
    onRetry: () => void
}

function ErrorState({ message, onRetry }: ErrorStateProps) {
    return (
        <div>
            <p>{message}</p>

            <button onClick={onRetry}>
                Retry
            </button>
        </div>
    )
}

export default ErrorState