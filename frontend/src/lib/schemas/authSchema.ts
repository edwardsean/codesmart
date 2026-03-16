import { z } from "zod"


const authSchema = z.object({
  username: z.string().min(3, "Username must be at least 3 characters"),
  email: z.string().email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  confirm_password: z.string().min(6, "Confirm password must be at least 6 characters"),
});

export const registerSchema = authSchema.refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirmPassword"], // error appears on confirmPassword field
});

export const loginSchema = authSchema.omit({username: true, confirm_password: true});


export type RegisterFormData = z.infer<typeof registerSchema>
export type LoginFormData = z.infer<typeof loginSchema>