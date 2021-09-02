import React, { useState } from 'react'
import { useHistory } from 'react-router-dom'
import styled from 'styled-components'
import { mutate } from 'swr'
import CONSTANTS from '../constants'

interface ILoginFormProps {}

export default function LoginForm(props: ILoginFormProps) {
    const [email, setEmail] = useState('')
    const [verifyStep, setVerifyStep] = useState(false)
    const [disabled, setDisabled] = useState(false)
    const history = useHistory()

    const handleTryLogin = async (event: any) => {
        event.preventDefault()
        setDisabled(true)

        const formData = new FormData(event.currentTarget)

        try {
            await fetch(CONSTANTS.API.AUTH.LOGIN, {
                method: 'POST',
                body: formData,
            })
        } catch (e) {
            console.error(e)
        }

        setEmail(event.target.email.value)
        event.target.reset()
        setVerifyStep(true)

        setDisabled(false)
    }

    const handleLogin = async (event: any) => {
        event.preventDefault()
        setDisabled(true)
        
        const formData = new FormData(event.currentTarget)
        formData.set("email", email)
        try {
            await fetch(CONSTANTS.API.AUTH.VERIFY, {
                method: 'POST',
                body: formData,
            })
        } catch (e) {
            console.error(e)
        }
        
        mutate(CONSTANTS.API.AUTH.USER)
        history.push(CONSTANTS.ROUTES.RIDES.MAIN)
    }

    return (
        <>
            {!verifyStep && (
                <StyledForm onSubmit={handleTryLogin}>
                    <FormItem>
                        <span>Email:</span>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            placeholder="Enter your email..."
                            required
                        />
                    </FormItem>

                    <input type="submit" value="Login" disabled={disabled} />
                </StyledForm>
            )}
            {verifyStep && (
                <StyledForm onSubmit={handleLogin}>
                    <FormItem>
                        <span>One Time Password</span>
                        <input
                            type="text"
                            inputMode="numeric"
                            name="otp"
                            placeholder="Enter your code here..."
                            autoComplete="one-time-code"
                            required
                        />
                    </FormItem>

                    <input type="submit" value="Login" disabled={disabled} />
                </StyledForm>
            )}
        </>
    )
}

const StyledForm = styled.form`
    width: 16rem;

    display: flex;
    flex-direction: column;
    padding: 1rem;
    gap: 1rem;

    border-radius: 1rem;
    border: 1px solid lightgrey;
`

const FormItem = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: start;
    gap: 0.25rem;
`
