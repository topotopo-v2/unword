type LoadingStateProps = {
    message: string
}


function LoadingState({ message }: LoadingStateProps) {
    return <p>{message}</p>
}

export default LoadingState