<script>
    import Button from '../components/Button.svelte';
    import Input from '../components/Input.svelte';
    import store_logo from '../resources/shops.svg';
    import FieldData from '../classes/FieldData';
    import { FieldStates } from '../classes/FieldData';
    import { server_name } from '../server_info';
    import { push } from 'svelte-spa-router';
    
    let username_field = new FieldData("login-user", /^[@A-Za-z_\d]+$/, "username");
    let password_field = new FieldData("passwd-user", /^[^\s]+$/, "password","password");
    
    let form_inputs = [
        username_field,
        password_field
    ];

    let session_key = -1;

    const awaitEnter = e => {
        if (e.key.toLowerCase() === "enter") {
            if (verifyLogin()) {
                doLogin();
            }
        }
    }

    const clearInputs = () => {
        form_inputs.forEach(fd => fd.clear());
    }

    const doLogin = () => {
        const username = username_field.getFieldValue();
        const passwd = password_field.getFieldValue();
        
        const request = new Request(`http://${server_name}/login?username=${username}&password=${passwd}`);
        fetch(request)
            .then(promise => promise.json())
            .then(response => {
                if( response.response !== undefined) {
                    session_key = response.response;
                    push(`/store/${session_key}`)
                } else if(response.error !== undefined) {
                    alert(response.error);
                } else {
                    alert("something odd just happend");
                }
                clearInputs();
            });

    }

    const checkForm = () => {
        if (verifyLogin()) {
            doLogin();
        }
    }

    const verifyLogin = () => {
        let is_login_ready = true;
        for(let fd of form_inputs) {
            if(fd.getFieldValue() === "") {
                console.log(`field ${fd.name} is empty`)
                fd.state = FieldStates.NORMAL;
                is_login_ready = false;
            } else if(!fd.isReady()) {
                console.log(`field ${fd.name} doesnt match`)
                fd.state = FieldStates.HAS_ERRORS;
                is_login_ready = false;
            } else {
                fd.state = FieldStates.READY;
            }
        }
        form_inputs = [...form_inputs]
        console.log(`login is ready == ${is_login_ready}`);
        return is_login_ready;
    }




</script>

<style>
    #login-container {
        width: 60%;
        height: 80vh;
        display: flex;
        flex-direction: column;
        margin: 10vh auto;
        justify-content: space-evenly;
        align-items: center;
    }

    #banner-container {
        width:5vw;
    }

    #login-fields {
        width: 60%;
    }

    .login-field-container {
        margin-bottom: 2.5vh;
    }

    #login-controls {
        display: flex;
        width: 40%;
        justify-content: space-around;
    }


</style>



<div id="login-container">
    <div id="banner-container">
        {@html store_logo}
    </div>
    <div id="login-fields">
        {#each form_inputs as fd}
            <div class="login-field-container">
                <Input
                    field_data={fd}
                    input_padding=".5vh 2.5vw"
                    onKeypressed={awaitEnter}
                    onBlur={verifyLogin}
                />
            </div>
        {/each}
    </div>
    <div id="login-controls">
        <Button width="5vw" label="SingUp" onClick={() => push("/sign-up")} isDimmed={true}/>
        <Button width="10vw" label="Login" onClick={checkForm}/>
    </div>
</div>