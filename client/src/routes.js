import Login from './pages/login/Login.svelte';
import SignUp from './pages/signUp/SignUp.svelte';
import Store from './pages/storeMain/StoreMainView.svelte';
import UserProfile from './pages/userProfile/UserProfile.svelte';


const routes = {
    "/": Login,
    "/sign-up": SignUp,
    "/store/:sk": Store,
    "/profile/:sk": UserProfile 
}

export { routes }