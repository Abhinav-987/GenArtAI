# AI Image Generation Web Application

This project is a web application that generates images based on user prompts using the Stable Diffusion model provided by [Replicate](https://replicate.com/stability-ai/stable-diffusion/examples). The application is built using a combination of modern technologies for both the backend and frontend.

## Technologies Used

### Backend
- **Go**: The programming language used to write the backend logic.
- **Chi**: A lightweight, idiomatic router for building HTTP services in Go.
- **Templ**: A powerful templating engine for Go, used to render HTML templates.
- **Supabase**: A backend-as-a-service providing database, authentication, and storage functionalities.
- **Stable Diffusion via Replicate**: An advanced AI model for generating images based on textual prompts, accessed through the Replicate API.

### Frontend
- **TailwindCSS**: A utility-first CSS framework for styling the application.
- **DaisyUI**: A component library that integrates with TailwindCSS to provide pre-built UI components.
- **htmx**: A library that allows you to add interactivity to your web application using HTML attributes.

## Features

- **Image Generation**: Users can input prompts, and the application will generate images using the Stable Diffusion model accessed via Replicate.
- **Dynamic UI**: The application uses htmx to make the interface dynamic and interactive without needing extensive JavaScript.
- **Responsive Design**: Styled with TailwindCSS and DaisyUI, the application is fully responsive and works on various devices.
- **Database Integration**: Supabase is used to handle database operations and user authentication.
