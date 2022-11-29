# Production build of tailwind
npx tailwindcss -o ../static/tailwind.css --minify

# Development
npx tailwindcss -i input.css -o ../static/tailwind.css --watch
