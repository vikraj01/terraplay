from manim import *

class KeeperInfra(Scene):
    def construct(self):
        self.show_title()
        terraform_logo = self.show_terraform_icon()
        self.show_aws_components(terraform_logo)
        self.show_latex_text()

    def show_title(self):
        title = Text("Keeper Infra Deployment", font_size=48)
        self.play(FadeIn(title, shift=UP), run_time=1.5)
        self.wait(1)
        self.play(title.animate.scale(0.7).to_edge(UP))

    def show_terraform_icon(self):
        terraform_logo = SVGMobject("assets/terraform.svg").scale(1).to_edge(LEFT, buff=1.5)
        self.play(FadeIn(terraform_logo, scale=1.1), run_time=1.5)
        self.wait(0.5)
        return terraform_logo

    def show_aws_components(self, terraform_logo):
        discord_bot, discord_label = self.add_component("assets/discordbot.svg", "Discord Bot (EC2 Server)", RIGHT * 4.5 + UP * 2.5)
        iam_icon, iam_label = self.add_component("assets/iam.svg", "IAM: For Permissions", RIGHT * 4.5 + UP * 1)
        dynamodb_icon, dynamo_label = self.add_component("assets/dynamodb.svg", "DynamoDB: For Locking", RIGHT * 4.5 + DOWN * 1)
        s3_icon, s3_label = self.add_component("assets/s3.svg", "S3: For State Management", RIGHT * 4.5 + DOWN * 2.5)

        self.animate_aws_component(discord_bot, discord_label)
        self.animate_aws_component(iam_icon, iam_label)
        self.animate_aws_component(dynamodb_icon, dynamo_label)
        self.animate_aws_component(s3_icon, s3_label)

    def add_component(self, svg_path, label_text, position):
        icon = SVGMobject(svg_path).scale(0.3).move_to(position)
        label = Text(label_text, font_size=12).next_to(icon, DOWN, buff=0.1)
        return icon, label

    def animate_aws_component(self, icon, label):
        self.play(FadeIn(icon, label), run_time=1.2)

    def show_latex_text(self):
        latex_text = MathTex(r"\text{Create them or Import them!}", font_size=48).to_edge(DOWN)
        self.play(Write(latex_text))
        self.wait(2)
