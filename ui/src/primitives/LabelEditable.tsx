import React, { createRef } from 'react'
import { TextField } from './form/Input'

type Props = {
    value?: string | null
    placeholder?: string
    loading?: boolean
    className: string
    /**
     * Called when the user confirms the edit (presses Enter or clicks the action).
     * Passes null when the label is empty.
     */
    onSave: (label: string | null) => void
}

type State = {
    editing: boolean
    inputValue: string
}

class LabelEditable extends React.Component<Props, State> {
    private inputRef = createRef<HTMLInputElement>()

    constructor(props: Props) {
        super(props)
        this.state = {
            editing: false,
            inputValue: props.value ?? '',
        }
    }

    componentDidUpdate(prevProps: Props, prevState: State) {
        // keep inputValue in sync when not editing
        console.log('LabelEditable componentDidUpdate', this.props.value, prevProps.value)
        if (!this.state.editing && prevProps.value !== this.props.value) {
            this.setState({ inputValue: this.props.value ?? '' })
        }

        // focus the input when entering edit mode
        if (this.state.editing && !prevState.editing) {
            if (this.inputRef.current) {
                this.inputRef.current.focus()
                this.inputRef.current.select()
            }
        }
    }

    enterEdit = () => {
        this.setState({ editing: true, inputValue: this.props.value ?? '' })
    }

    exitAndReset = () => {
        this.setState({ editing: false, inputValue: this.props.value ?? '' })
    }

    save = () => {
        const val = this.state.inputValue.trim()
        this.props.onSave(val ? val : null)
        this.setState({ editing: false })
    }

    handleKeyDown: React.KeyboardEventHandler<HTMLInputElement> = e => {
        if (e.key === 'Escape') {
            this.exitAndReset()
            return
        }
        if (e.key === 'Enter') {
            this.save()
            return
        }
    }

    render() {
        const { className, placeholder, loading } = this.props
        const { editing, inputValue } = this.state

        if (!editing) {
            return (
                <h1
                    className={className}
                    onClick={this.enterEdit}
                    role="button"
                    tabIndex={0}
                    onKeyDown={e => {
                        if (e.key === 'Enter' || e.key === ' ') {
                            this.enterEdit()
                        }
                    }}
                >
                    {this.props.value ?? placeholder ?? ''}
                </h1>
            )
        }

        return (
            <TextField
                loading={loading}
                ref={this.inputRef}
                placeholder={placeholder ?? ''}
                action={() => {
                    this.save()
                }}
                value={inputValue}
                onKeyDown={this.handleKeyDown}
                onChange={e => this.setState({ inputValue: e.target.value })}
                onBlur={() => {
                    this.exitAndReset()
                }}
            />
        )
    }
}

export default LabelEditable