



export async function login(email: string, password: string) {
    try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_GOLANG_API_URL}/login`, {
            method: "POST",
            headers : {
                "Content-Type" : "application/json",
            },
            body : JSON.stringify({email, password})
        });

        if (!response.ok) {
            throw new Error()
        }

        return response.json() //returns {token : ""}
    } catch (err) {
        throw err
    }
}

// export async function register(email: string, password: string) {
//     try {
//         const response = await fetch("http://localhost:8080/api/v1/login", {
//             method: "POST",
//             headers : {
//                 "Content-Type" : "application/json",
//             },
//             body : JSON.stringify({email, password})
//         });

//         if (!response.ok) {
            
//         }
//     } catch (err) {
//         throw err
//     }
// }

