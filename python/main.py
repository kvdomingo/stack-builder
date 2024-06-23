import typer
from rich import print
import inquirer

app = typer.Typer()


@app.command()
def main():
    questions = [
        inquirer.Confirm(
            "is_typescript", message="Do you want to use TypeScript?", default=True
        ),
        inquirer.List(
            "framework",
            message="Choose a framework",
            choices=[
                "Express",
                "Hono",
                "Nest",
                "Next",
                "Nuxt",
                "SvelteKit",
                "SolidStart",
                "Vite + React",
                "Vite + Vue",
                "Vite + Svelte",
                "Vite + Solid",
            ],
            carousel=True,
        ),
        inquirer.List(
            "ui_framework",
            message="Choose a UI framework",
            choices=[
                "None",
                "shadcn/ui",
                "Tailwind UI",
                "Daisy UI",
                "Material UI",
                "Chakra",
                "Semantic UI",
            ],
            carousel=True,
        ),
        inquirer.List(
            "css_framework",
            message="Choose a CSS framework",
            choices=[
                "None",
                "TailwindCSS",
                "Emotion",
                "StyledCSS",
                "Bootstrap",
            ],
            carousel=True,
        ),
        inquirer.List(
            "database",
            message="Choose a relational database management system",
            choices=["None", "PostgreSQL", "MySQL", "SQLite"],
            carousel=True,
        ),
        inquirer.List(
            "orm",
            message="Choose an object-relational mapper",
            choices=["None", "Prisma", "Drizzle"],
            carousel=True,
        ),
        inquirer.List(
            "cloud_platform",
            message="Where will you be deploying to?",
            choices=[
                "None",
                "GCP Cloud Run",
                "GCP Cloud Functions",
                "GCP App Engine",
                "AWS Lambda",
                "AWS Elastic Container Service",
                "Azure Container Service",
                "Vercel",
            ],
            carousel=True,
        ),
    ]

    answers = inquirer.prompt(questions)

    print(answers)


if __name__ == "__main__":
    typer.run(main)
