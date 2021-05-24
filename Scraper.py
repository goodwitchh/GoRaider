import asyncio
import discord, json
from discord.ext import commands

with open("LoginInfo.json") as file:
    data = json.load(file)
    file.close()
    
client = discord.Client() 
client = commands.Bot(".", self_bot=data["Bot"], chunk_guilds_at_startup=False)
client.remove_command('help')

@client.event
async def on_connect():
    guild = client.get_guild(int(data['GuildID']))
    for channel in guild.text_channels:
        if channel.permissions_for(guild.me).send_messages:
            await channel.send(".s")
        break

@client.command(aliases=['s'])
async def scrape(ctx):
    scrapedamout = 0
    await ctx.message.delete()
    memberslist = []
    for member in ctx.guild.members:
        if member.id == client.user.id:
            continue
        memberslist.append(str(member.id))
        scrapedamout += 1
    membersjson = json.dumps(memberslist, separators=(',', ':'))
    print(membersjson)
    await client.close()
client.run(data["Token"]) # idk why you would use this for a bot acc but /shrug
