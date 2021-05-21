<script>
    export let label = "button";
    export let isEnabled = true;
    export let isDimmed = false;
    export let onClick = () => {};
    export let isDanger = false;
    export let width = "100%";

    const composeClassName = () => {
        let class_name = "control-btn"
        // Red class
        if (isDanger) {
            class_name += " red";
        }

        //Enable State
        class_name += isEnabled ? " enable" : " disable";


        //Dimmed class
        if (isDimmed) {
            class_name += " dimmed";
        }
        
        return class_name;
    }

</script>

<style>

    @keyframes Enable {
        0% {filter: grayscale(90%);}
        25% {filter: grayscale(60%);}
        50% {filter: grayscale(40%);}
        75% {filter: grayscale(20%);}
        100% {filter: grayscale(0%);}
    }

    .control-btn {
        font-family: var(--main-font);
        cursor: pointer;
        background-color: var(--theme-color);
        color: var(--dark-color);
        text-align: center;
        padding: 0 5%;
        border: none;
        border-radius: 15px;
        transition: all .4s ease-in-out;
    }

    .control-btn.red {
        background-color: var(--danger) !important;
    }

    .control-btn.enable {
        animation-name: Enable;
        animation-duration: 1s;
        animation-iteration-count: 1;
    }

    .control-btn.disable {
        filter: grayscale(90%);
    }

    .control-btn.dimmed {   
        background: none !important;
        color: var(--theme-color);
        border: 1px solid var(--theme-color);
    }

    .control-btn.dimmed:hover {
        color: white;
        background-color: var(--theme-color) !important;
    }

</style>

<div style="width: {width};" on:click={isEnabled ?  onClick : ()=>{}}
    class={composeClassName()}>
        {label}
</div>