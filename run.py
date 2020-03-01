from twitchio.ext import commands
from utils.config import load_config
import datetime
import os

absolute_path = os.path.dirname(os.path.abspath(__file__))
config = load_config(f'{absolute_path}/config.yml')


class Bot(commands.Bot):
    cooldowns = {}

    def __init__(self, irc_token, nick, initial_channels, prefix='!'):
        super().__init__(irc_token=irc_token, nick=nick, initial_channels=initial_channels, prefix=prefix)

    def cooldown(self, cmd) -> bool:
        # check if cooldown (30 seconds) is reached for given command
        if (cmd not in self.cooldowns.keys() or
                self.cooldowns[cmd] + datetime.timedelta(seconds=30) < datetime.datetime.now()):
            self.cooldowns[cmd] = datetime.datetime.now()
            return False

        return True

    async def event_ready(self):
        print(f'Connected to Twitch chat - {self.nick}')

    async def event_message(self, message):
        await self.handle_commands(message)

    @commands.command(name='vscode', aliases=['editor'])
    async def vscode(self, ctx):
        if self.cooldown('vscode'):
            return
        await ctx.send('Benutzt wird vscode von Microsoft - https://code.visualstudio.com/')

    @commands.command(name='extensions', aliases=['ext'])
    async def extensions(self, ctx):
        if self.cooldown('extensions'):
            return
        await ctx.send('Die für mich wichtigsten Extensions (für vscode) sind: Remote Containers, Settings Sync, TODO Highlight, GitLens, markdownlint, Prettier, YAML, Todo+, Todo Tree, TODO Highlight')

    @commands.command(name='repository', aliases=['repo', 'repositories', 'repos'])
    async def repository(self, ctx):
        if self.cooldown('repository'):
            return
        await ctx.send('Alle verwendeten Repositories findest du unter https://github.com/curi0s - Das Repository zum Thema "Tech Streams" findest du hier https://github.com/curi0s/stream')

    @commands.command(name='hackintosh', aliases=['hack', 'mac', 'macos'])
    async def hackintosh(self, ctx):
        if self.cooldown('hackintosh'):
            return
        await ctx.send('Der Hackintosh ist ein PC auf dem MacOS installiert ist. Die verbauten Komponenten findest du hier https://github.com/curi0s/stream#computer-hackintosh')

    @commands.command(name='job', aliases=['beruf'])
    async def job(self, ctx):
        if self.cooldown('job'):
            return
        await ctx.send('Fabian hat mit ca. 12 Jahren begonnen, sich fürs Programmieren zu interessieren und hat dadurch mit HTML, CSS und PHP angefangen. Später hat er dann eine Ausbildung zum Fachinformatiker Systemintegration absolviert und arbeitet inzwischen seit mehr als 10 Jahren als Systemengineer')

    @commands.command(name='theme')
    async def theme(self, ctx):
        if self.cooldown('theme'):
            return
        await ctx.send('Das Theme für vscode ist JetJet-Alternate-Gray - https://marketplace.visualstudio.com/items?itemName=JohnyGeorges.jetjet-theme')

    @commands.command(name='today', aliases=['heute'])
    async def today(self, ctx):
        if self.cooldown('today') or not os.path.isfile(f'{absolute_path}/today.txt'):
            return

        with open(f'{absolute_path}/today.txt', 'r') as file:
            content = file.read()
        await ctx.send(content)

    @commands.command(name='social', aliases=['twitter', 'github', 'git', 'discord', 'insta', 'instagram'])
    async def social(self, ctx):
        if self.cooldown('social'):
            return
        await ctx.send('Twitter: https://www.twitter.com/curi0sDE - Discord: https://discord.gg/curi0sDE - Instagram: https://www.instagram.com/curi0sDE')

    @commands.command(name='font')
    async def font(self, ctx):
        if self.cooldown('font'):
            return
        await ctx.send('Als font wird Fira Code V2 verwendet: https://github.com/tonsky/FiraCode')


bot = Bot(
    config['irc_token'],
    config['nick'],
    config['initial_channels'],
    config['prefix'])
bot.run()
