package minecraft

const deathRegexps = `
\S+ was shot by \S+
\S+ was shot by \S+ using \S+
\S+ was pummeled by \S+
\S+ was pummeled by \S+ using \S+
\S+ was pricked to death
\S+ walked into a cactus whilst trying to escape \S+
\S+ drowned
\S+ drowned whilst trying to escape \S+
\S+ experienced kinetic energy
\S+ experienced kinetic energy whilst trying to escape \S+
\S+ blew up
\S+ was blown up by \S+
\S+ was blown up by \S+ using \S+
\S+ was killed by \S+
\S+ hit the ground too hard
\S+ hit the ground too hard whilst trying to escape \S+
\S+ fell from a high place
\S+ fell off a ladder
\S+ fell off some vines
\S+ fell off some weeping vines
\S+ fell off some twisting vines
\S+ fell off scaffolding
\S+ fell while climbing
\S+ was impaled on a stalagmite
\S+ was impaled on a stalagmite whilst fighting \S+
\S+ was squashed by a falling anvil
\S+ was squashed by a falling anvil whilst fighting \S+
\S+ was squashed by a falling block
\S+ was squashed by a falling block whilst fighting \S+
\S+ was skewered by a falling stalactite
\S+ was skewered by a falling stalactite whilst fighting \S+
\S+ went up in flames
\S+ walked into fire whilst fighting \S+
\S+ burned to death
\S+ was burnt to a crisp whilst fighting \S+
\S+ went off with a bang
\S+ went off with a bang due to a firework fired from \S+ by \S+
\S+ tried to swim in lava
\S+ tried to swim in lava to escape \S+
\S+ was struck by lightning
\S+ was struck by lightning whilst fighting \S+
\S+ discovered the floor was lava
\S+ walked into danger zone due to \S+
\S+ was killed by magic
\S+ was killed by magic whilst trying to escape \S+
\S+ was killed by \S+ using magic
\S+ was killed by \S+ using \S+
\S+ froze to death
\S+ was frozen to death by \S+
\S+ was slain by \S+
\S+ was slain by \S+ using \S+
\S+ was fireballed by \S+
\S+ was fireballed by \S+ using \S+
\S+ was stung to death
\S+ was shot by a skull from \S+
\S+ starved to death
\S+ starved to death whilst fighting \S+
\S+ suffocated in a wall
\S+ suffocated in a wall whilst fighting \S+
\S+ was squished too much
\S+ was squashed by \S+
\S+ was poked to death by a sweet berry bush
\S+ was poked to death by a sweet berry bush whilst trying to escape \S+
\S+ was killed trying to hurt \S+
\S+ was killed by \S+ trying to hurt \S+
\S+ was impaled by \S+
\S+ was impaled by \S+ with \S+
\S+ fell out of the world
\S+ didnt want to live in the same world as \S+
\S+ withered away
\S+ withered away whilst fighting \S+
\S+ died from dehydration
\S+ died from dehydration whilst trying to escape \S+
\S+ died
\S+ died because of \S+
\S+ was roasted in dragon breath
\S+ was roasted in dragon breath by \S+
\S+ was doomed to fall
\S+ was doomed to fall by \S+
\S+ was doomed to fall by \S+ using \S+
\S+ fell too far and was finished by \S+
\S+ fell too far and was finished by \S+ using \S+
\S+ was stung to death by \S+
\S+ went off with a bang whilst fighting \S+
\S+ was killed by even more magic
\S+ was too soft for this world
\S+ was too soft for this world (\S+ helped)
\S+ fell (off a ladder|off some vines|out of the water|from a high place)
\S+ was (shot|pummeled|blown|fireballed|knocked) (off a ladder|off some vines|out of the water|from a high place)
\S+ was (shot|pummeled|blown|fireballed|knocked) (off a ladder|off some vines|out of the water|from a high place) by \S+
\S+ into a pool of lava
\S+ into a patch of fire
\S+ into a patch of cacti
\S+ got finished off by \S+
\S+ got finished off by \S+ using \S+
\S+ got finished off using \S+
\S+ fell off a ladder
\S+ fell off some vines
\S+ fell out of the water
\S+ fell from a high place
\S+ was shot off a ladder
\S+ was shot off some vines
\S+ was shot out of the water
\S+ was shot from a high place
\S+ was pummeled off a ladder
\S+ was pummeled off some vines
\S+ was pummeled out of the water
\S+ was pummeled from a high place
\S+ was blown off a ladder
\S+ was blown off some vines
\S+ was blown out of the water
\S+ was blown from a high place
\S+ was fireballed off a ladder
\S+ was fireballed off some vines
\S+ was fireballed out of the water
\S+ was fireballed from a high place
\S+ was knocked off a ladder
\S+ was knocked off some vines
\S+ was knocked out of the water
\S+ was knocked from a high place
\S+ was shot off a ladder by \S+
\S+ was shot off some vines by \S+
\S+ was shot out of the water by \S+
\S+ was shot from a high place by \S+
\S+ was pummeled off a ladder by \S+
\S+ was pummeled off some vines by \S+
\S+ was pummeled out of the water by \S+
\S+ was pummeled from a high place by \S+
\S+ was blown off a ladder by \S+
\S+ was blown off some vines by \S+
\S+ was blown out of the water by \S+
\S+ was blown from a high place by \S+
\S+ was fireballed off a ladder by \S+
\S+ was fireballed off some vines by \S+
\S+ was fireballed out of the water by \S+
\S+ was fireballed from a high place by \S+
\S+ was knocked off a ladder by \S+
\S+ was knocked off some vines by \S+
\S+ was knocked out of the water by \S+
\S+ was knocked from a high place by \S+
\S+ fell off a ladder and into a pool of lava
\S+ fell off some vines and into a pool of lava
\S+ fell out of the water and into a pool of lava
\S+ fell from a high place and into a pool of lava
\S+ was shot off a ladder and into a pool of lava
\S+ was shot off some vines and into a pool of lava
\S+ was shot out of the water and into a pool of lava
\S+ was shot from a high place and into a pool of lava
\S+ was pummeled off a ladder and into a pool of lava
\S+ was pummeled off some vines and into a pool of lava
\S+ was pummeled out of the water and into a pool of lava
\S+ was pummeled from a high place and into a pool of lava
\S+ was blown off a ladder and into a pool of lava
\S+ was blown off some vines and into a pool of lava
\S+ was blown out of the water and into a pool of lava
\S+ was blown from a high place and into a pool of lava
\S+ was fireballed off a ladder and into a pool of lava
\S+ was fireballed off some vines and into a pool of lava
\S+ was fireballed out of the water and into a pool of lava
\S+ was fireballed from a high place and into a pool of lava
\S+ was knocked off a ladder and into a pool of lava
\S+ was knocked off some vines and into a pool of lava
\S+ was knocked out of the water and into a pool of lava
\S+ was knocked from a high place and into a pool of lava
\S+ was shot off a ladder by \S+ and into a pool of lava
\S+ was shot off some vines by \S+ and into a pool of lava
\S+ was shot out of the water by \S+ and into a pool of lava
\S+ was shot from a high place by \S+ and into a pool of lava
\S+ was pummeled off a ladder by \S+ and into a pool of lava
\S+ was pummeled off some vines by \S+ and into a pool of lava
\S+ was pummeled out of the water by \S+ and into a pool of lava
\S+ was pummeled from a high place by \S+ and into a pool of lava
\S+ was blown off a ladder by \S+ and into a pool of lava
\S+ was blown off some vines by \S+ and into a pool of lava
\S+ was blown out of the water by \S+ and into a pool of lava
\S+ was blown from a high place by \S+ and into a pool of lava
\S+ was fireballed off a ladder by \S+ and into a pool of lava
\S+ was fireballed off some vines by \S+ and into a pool of lava
\S+ was fireballed out of the water by \S+ and into a pool of lava
\S+ was fireballed from a high place by \S+ and into a pool of lava
\S+ was knocked off a ladder by \S+ and into a pool of lava
\S+ was knocked off some vines by \S+ and into a pool of lava
\S+ was knocked out of the water by \S+ and into a pool of lava
\S+ was knocked from a high place by \S+ and into a pool of lava
\S+ fell off a ladder and fell out of the world
\S+ fell off some vines and fell out of the world
\S+ fell out of the water and fell out of the world
\S+ fell from a high place and fell out of the world
\S+ was shot off a ladder and fell out of the world
\S+ was shot off some vines and fell out of the world
\S+ was shot out of the water and fell out of the world
\S+ was shot from a high place and fell out of the world
\S+ was pummeled off a ladder and fell out of the world
\S+ was pummeled off some vines and fell out of the world
\S+ was pummeled out of the water and fell out of the world
\S+ was pummeled from a high place and fell out of the world
\S+ was blown off a ladder and fell out of the world
\S+ was blown off some vines and fell out of the world
\S+ was blown out of the water and fell out of the world
\S+ was blown from a high place and fell out of the world
\S+ was fireballed off a ladder and fell out of the world
\S+ was fireballed off some vines and fell out of the world
\S+ was fireballed out of the water and fell out of the world
\S+ was fireballed from a high place and fell out of the world
\S+ was knocked off a ladder and fell out of the world
\S+ was knocked off some vines and fell out of the world
\S+ was knocked out of the water and fell out of the world
\S+ was knocked from a high place and fell out of the world
\S+ was shot off a ladder by \S+ and fell out of the world
\S+ was shot off some vines by \S+ and fell out of the world
\S+ was shot out of the water by \S+ and fell out of the world
\S+ was shot from a high place by \S+ and fell out of the world
\S+ was pummeled off a ladder by \S+ and fell out of the world
\S+ was pummeled off some vines by \S+ and fell out of the world
\S+ was pummeled out of the water by \S+ and fell out of the world
\S+ was pummeled from a high place by \S+ and fell out of the world
\S+ was blown off a ladder by \S+ and fell out of the world
\S+ was blown off some vines by \S+ and fell out of the world
\S+ was blown out of the water by \S+ and fell out of the world
\S+ was blown from a high place by \S+ and fell out of the world
\S+ was fireballed off a ladder by \S+ and fell out of the world
\S+ was fireballed off some vines by \S+ and fell out of the world
\S+ was fireballed out of the water by \S+ and fell out of the world
\S+ was fireballed from a high place by \S+ and fell out of the world
\S+ was knocked off a ladder by \S+ and fell out of the world
\S+ was knocked off some vines by \S+ and fell out of the world
\S+ was knocked out of the water by \S+ and fell out of the world
\S+ was knocked from a high place by \S+ and fell out of the world
\S+ fell off a ladder and into a patch of fire
\S+ fell off some vines and into a patch of fire
\S+ fell out of the water and into a patch of fire
\S+ fell from a high place and into a patch of fire
\S+ was shot off a ladder and into a patch of fire
\S+ was shot off some vines and into a patch of fire
\S+ was shot out of the water and into a patch of fire
\S+ was shot from a high place and into a patch of fire
\S+ was pummeled off a ladder and into a patch of fire
\S+ was pummeled off some vines and into a patch of fire
\S+ was pummeled out of the water and into a patch of fire
\S+ was pummeled from a high place and into a patch of fire
\S+ was blown off a ladder and into a patch of fire
\S+ was blown off some vines and into a patch of fire
\S+ was blown out of the water and into a patch of fire
\S+ was blown from a high place and into a patch of fire
\S+ was fireballed off a ladder and into a patch of fire
\S+ was fireballed off some vines and into a patch of fire
\S+ was fireballed out of the water and into a patch of fire
\S+ was fireballed from a high place and into a patch of fire
\S+ was knocked off a ladder and into a patch of fire
\S+ was knocked off some vines and into a patch of fire
\S+ was knocked out of the water and into a patch of fire
\S+ was knocked from a high place and into a patch of fire
\S+ was shot off a ladder by \S+ and into a patch of fire
\S+ was shot off some vines by \S+ and into a patch of fire
\S+ was shot out of the water by \S+ and into a patch of fire
\S+ was shot from a high place by \S+ and into a patch of fire
\S+ was pummeled off a ladder by \S+ and into a patch of fire
\S+ was pummeled off some vines by \S+ and into a patch of fire
\S+ was pummeled out of the water by \S+ and into a patch of fire
\S+ was pummeled from a high place by \S+ and into a patch of fire
\S+ was blown off a ladder by \S+ and into a patch of fire
\S+ was blown off some vines by \S+ and into a patch of fire
\S+ was blown out of the water by \S+ and into a patch of fire
\S+ was blown from a high place by \S+ and into a patch of fire
\S+ was fireballed off a ladder by \S+ and into a patch of fire
\S+ was fireballed off some vines by \S+ and into a patch of fire
\S+ was fireballed out of the water by \S+ and into a patch of fire
\S+ was fireballed from a high place by \S+ and into a patch of fire
\S+ was knocked off a ladder by \S+ and into a patch of fire
\S+ was knocked off some vines by \S+ and into a patch of fire
\S+ was knocked out of the water by \S+ and into a patch of fire
\S+ was knocked from a high place by \S+ and into a patch of fire
\S+ fell off a ladder and into a patch of cacti
\S+ fell off some vines and into a patch of cacti
\S+ fell out of the water and into a patch of cacti
\S+ fell from a high place and into a patch of cacti
\S+ was shot off a ladder and into a patch of cacti
\S+ was shot off some vines and into a patch of cacti
\S+ was shot out of the water and into a patch of cacti
\S+ was shot from a high place and into a patch of cacti
\S+ was pummeled off a ladder and into a patch of cacti
\S+ was pummeled off some vines and into a patch of cacti
\S+ was pummeled out of the water and into a patch of cacti
\S+ was pummeled from a high place and into a patch of cacti
\S+ was blown off a ladder and into a patch of cacti
\S+ was blown off some vines and into a patch of cacti
\S+ was blown out of the water and into a patch of cacti
\S+ was blown from a high place and into a patch of cacti
\S+ was fireballed off a ladder and into a patch of cacti
\S+ was fireballed off some vines and into a patch of cacti
\S+ was fireballed out of the water and into a patch of cacti
\S+ was fireballed from a high place and into a patch of cacti
\S+ was knocked off a ladder and into a patch of cacti
\S+ was knocked off some vines and into a patch of cacti
\S+ was knocked out of the water and into a patch of cacti
\S+ was knocked from a high place and into a patch of cacti
\S+ was shot off a ladder by \S+ and into a patch of cacti
\S+ was shot off some vines by \S+ and into a patch of cacti
\S+ was shot out of the water by \S+ and into a patch of cacti
\S+ was shot from a high place by \S+ and into a patch of cacti
\S+ was pummeled off a ladder by \S+ and into a patch of cacti
\S+ was pummeled off some vines by \S+ and into a patch of cacti
\S+ was pummeled out of the water by \S+ and into a patch of cacti
\S+ was pummeled from a high place by \S+ and into a patch of cacti
\S+ was blown off a ladder by \S+ and into a patch of cacti
\S+ was blown off some vines by \S+ and into a patch of cacti
\S+ was blown out of the water by \S+ and into a patch of cacti
\S+ was blown from a high place by \S+ and into a patch of cacti
\S+ was fireballed off a ladder by \S+ and into a patch of cacti
\S+ was fireballed off some vines by \S+ and into a patch of cacti
\S+ was fireballed out of the water by \S+ and into a patch of cacti
\S+ was fireballed from a high place by \S+ and into a patch of cacti
\S+ was knocked off a ladder by \S+ and into a patch of cacti
\S+ was knocked off some vines by \S+ and into a patch of cacti
\S+ was knocked out of the water by \S+ and into a patch of cacti
\S+ was knocked from a high place by \S+ and into a patch of cacti
\S+ fell off a ladder and got finished off by \S+
\S+ fell off some vines and got finished off by \S+
\S+ fell out of the water and got finished off by \S+
\S+ fell from a high place and got finished off by \S+
\S+ was shot off a ladder and got finished off by \S+
\S+ was shot off some vines and got finished off by \S+
\S+ was shot out of the water and got finished off by \S+
\S+ was shot from a high place and got finished off by \S+
\S+ was pummeled off a ladder and got finished off by \S+
\S+ was pummeled off some vines and got finished off by \S+
\S+ was pummeled out of the water and got finished off by \S+
\S+ was pummeled from a high place and got finished off by \S+
\S+ was blown off a ladder and got finished off by \S+
\S+ was blown off some vines and got finished off by \S+
\S+ was blown out of the water and got finished off by \S+
\S+ was blown from a high place and got finished off by \S+
\S+ was fireballed off a ladder and got finished off by \S+
\S+ was fireballed off some vines and got finished off by \S+
\S+ was fireballed out of the water and got finished off by \S+
\S+ was fireballed from a high place and got finished off by \S+
\S+ was knocked off a ladder and got finished off by \S+
\S+ was knocked off some vines and got finished off by \S+
\S+ was knocked out of the water and got finished off by \S+
\S+ was knocked from a high place and got finished off by \S+
\S+ was shot off a ladder by \S+ and got finished off by \S+
\S+ was shot off some vines by \S+ and got finished off by \S+
\S+ was shot out of the water by \S+ and got finished off by \S+
\S+ was shot from a high place by \S+ and got finished off by \S+
\S+ was pummeled off a ladder by \S+ and got finished off by \S+
\S+ was pummeled off some vines by \S+ and got finished off by \S+
\S+ was pummeled out of the water by \S+ and got finished off by \S+
\S+ was pummeled from a high place by \S+ and got finished off by \S+
\S+ was blown off a ladder by \S+ and got finished off by \S+
\S+ was blown off some vines by \S+ and got finished off by \S+
\S+ was blown out of the water by \S+ and got finished off by \S+
\S+ was blown from a high place by \S+ and got finished off by \S+
\S+ was fireballed off a ladder by \S+ and got finished off by \S+
\S+ was fireballed off some vines by \S+ and got finished off by \S+
\S+ was fireballed out of the water by \S+ and got finished off by \S+
\S+ was fireballed from a high place by \S+ and got finished off by \S+
\S+ was knocked off a ladder by \S+ and got finished off by \S+
\S+ was knocked off some vines by \S+ and got finished off by \S+
\S+ was knocked out of the water by \S+ and got finished off by \S+
\S+ was knocked from a high place by \S+ and got finished off by \S+
\S+ fell off a ladder and got finished off by \S+ using \S+
\S+ fell off some vines and got finished off by \S+ using \S+
\S+ fell out of the water and got finished off by \S+ using \S+
\S+ fell from a high place and got finished off by \S+ using \S+
\S+ was shot off a ladder and got finished off by \S+ using \S+
\S+ was shot off some vines and got finished off by \S+ using \S+
\S+ was shot out of the water and got finished off by \S+ using \S+
\S+ was shot from a high place and got finished off by \S+ using \S+
\S+ was pummeled off a ladder and got finished off by \S+ using \S+
\S+ was pummeled off some vines and got finished off by \S+ using \S+
\S+ was pummeled out of the water and got finished off by \S+ using \S+
\S+ was pummeled from a high place and got finished off by \S+ using \S+
\S+ was blown off a ladder and got finished off by \S+ using \S+
\S+ was blown off some vines and got finished off by \S+ using \S+
\S+ was blown out of the water and got finished off by \S+ using \S+
\S+ was blown from a high place and got finished off by \S+ using \S+
\S+ was fireballed off a ladder and got finished off by \S+ using \S+
\S+ was fireballed off some vines and got finished off by \S+ using \S+
\S+ was fireballed out of the water and got finished off by \S+ using \S+
\S+ was fireballed from a high place and got finished off by \S+ using \S+
\S+ was knocked off a ladder and got finished off by \S+ using \S+
\S+ was knocked off some vines and got finished off by \S+ using \S+
\S+ was knocked out of the water and got finished off by \S+ using \S+
\S+ was knocked from a high place and got finished off by \S+ using \S+
\S+ was shot off a ladder by \S+ and got finished off by \S+ using \S+
\S+ was shot off some vines by \S+ and got finished off by \S+ using \S+
\S+ was shot out of the water by \S+ and got finished off by \S+ using \S+
\S+ was shot from a high place by \S+ and got finished off by \S+ using \S+
\S+ was pummeled off a ladder by \S+ and got finished off by \S+ using \S+
\S+ was pummeled off some vines by \S+ and got finished off by \S+ using \S+
\S+ was pummeled out of the water by \S+ and got finished off by \S+ using \S+
\S+ was pummeled from a high place by \S+ and got finished off by \S+ using \S+
\S+ was blown off a ladder by \S+ and got finished off by \S+ using \S+
\S+ was blown off some vines by \S+ and got finished off by \S+ using \S+
\S+ was blown out of the water by \S+ and got finished off by \S+ using \S+
\S+ was blown from a high place by \S+ and got finished off by \S+ using \S+
\S+ was fireballed off a ladder by \S+ and got finished off by \S+ using \S+
\S+ was fireballed off some vines by \S+ and got finished off by \S+ using \S+
\S+ was fireballed out of the water by \S+ and got finished off by \S+ using \S+
\S+ was fireballed from a high place by \S+ and got finished off by \S+ using \S+
\S+ was knocked off a ladder by \S+ and got finished off by \S+ using \S+
\S+ was knocked off some vines by \S+ and got finished off by \S+ using \S+
\S+ was knocked out of the water by \S+ and got finished off by \S+ using \S+
\S+ was knocked from a high place by \S+ and got finished off by \S+ using \S+
\S+ fell off a ladder and got finished off using \S+
\S+ fell off some vines and got finished off using \S+
\S+ fell out of the water and got finished off using \S+
\S+ fell from a high place and got finished off using \S+
\S+ was shot off a ladder and got finished off using \S+
\S+ was shot off some vines and got finished off using \S+
\S+ was shot out of the water and got finished off using \S+
\S+ was shot from a high place and got finished off using \S+
\S+ was pummeled off a ladder and got finished off using \S+
\S+ was pummeled off some vines and got finished off using \S+
\S+ was pummeled out of the water and got finished off using \S+
\S+ was pummeled from a high place and got finished off using \S+
\S+ was blown off a ladder and got finished off using \S+
\S+ was blown off some vines and got finished off using \S+
\S+ was blown out of the water and got finished off using \S+
\S+ was blown from a high place and got finished off using \S+
\S+ was fireballed off a ladder and got finished off using \S+
\S+ was fireballed off some vines and got finished off using \S+
\S+ was fireballed out of the water and got finished off using \S+
\S+ was fireballed from a high place and got finished off using \S+
\S+ was knocked off a ladder and got finished off using \S+
\S+ was knocked off some vines and got finished off using \S+
\S+ was knocked out of the water and got finished off using \S+
\S+ was knocked from a high place and got finished off using \S+
\S+ was shot off a ladder by \S+ and got finished off using \S+
\S+ was shot off some vines by \S+ and got finished off using \S+
\S+ was shot out of the water by \S+ and got finished off using \S+
\S+ was shot from a high place by \S+ and got finished off using \S+
\S+ was pummeled off a ladder by \S+ and got finished off using \S+
\S+ was pummeled off some vines by \S+ and got finished off using \S+
\S+ was pummeled out of the water by \S+ and got finished off using \S+
\S+ was pummeled from a high place by \S+ and got finished off using \S+
\S+ was blown off a ladder by \S+ and got finished off using \S+
\S+ was blown off some vines by \S+ and got finished off using \S+
\S+ was blown out of the water by \S+ and got finished off using \S+
\S+ was blown from a high place by \S+ and got finished off using \S+
\S+ was fireballed off a ladder by \S+ and got finished off using \S+
\S+ was fireballed off some vines by \S+ and got finished off using \S+
\S+ was fireballed out of the water by \S+ and got finished off using \S+
\S+ was fireballed from a high place by \S+ and got finished off using \S+
\S+ was knocked off a ladder by \S+ and got finished off using \S+
\S+ was knocked off some vines by \S+ and got finished off using \S+
\S+ was knocked out of the water by \S+ and got finished off using \S+
\S+ was knocked from a high place by \S+ and got finished off using \S+
\S+ was slain by Arrow
\S+ was shot by \S+
\S+ was pricked to death
\S+ drowned
\S+ experienced kinetic energy
\S+ blew up
\S+ was blown up by Block of TNT
\S+ was blown up by \S+
\S+ hit the ground too hard
\S+ fell from a high place
\S+ was squashed by a falling anvil
\S+ was squashed by a falling block
\S+ went up in flames
\S+ burned to death
\S+ went off with a bang
\S+ tried to swim in lava
\S+ was struck by lightning
\S+ discovered floor was lava
\S+ was killed by magic
\S+ was killed by \S+ using magic
\S+ was slain by \S+
\S+ was slain by Small Fireball
\S+ starved to death
\S+ suffocated in a wall
\S+ was killed trying to hurt \S+
\S+ was impaled to death by \S+
\S+ fell out of the world
\S+ withered away
\S+ died
\S+ was fireballed by \S+
\S+ was sniped by \S+
\S+ was spitballed by \S+
\S+ froze to death
\S+ was skewered by a falling stalactite
\S+ was impaled on a stalagmite
\S+ was shot by \S+ using \S+
\S+ walked into a cactus whilst trying to escape \S+
\S+ drowned whilst trying to escape \S+
\S+ was fireballed by \S+ using \S+
\S+ was killed by \S+ using \S+
\S+ walked into fire whilst fighting \S+
\S+ tried to swim in lava to escape \S+
\S+ walked on danger zone due to \S+
\S+ was burnt to a crisp whilst fighting \S+
\S+ was slain by \S+ using \S+
\S+ was pummeled by \S+
\S+ was pummeled by \S+ using \S+
\S+ fell off a ladder
\S+ fell off some vines
\S+ fell out of the water
\S+ was doomed to fall
\S+ was doomed to fall by \S+
\S+ was doomed to fall by \S+ using \S+
\S+ fell too far and was finished by \S+
\S+ fell too far and was finished by \S+ using \S+
\S+ went up in flames
\S+ burned to death
\S+ tried to swim in lava
\S+ suffocated in a wall
\S+ drowned
\S+ starved to death
\S+ was pricked to death
\S+ hit the ground too hard
\S+ fell out of the world
\S+ died
\S+ blew up
\S+ was killed by magic
\S+ was slain by \S+
\S+ was slain by \S+
\S+ was shot by \S+
\S+ was fireballed by \S+
\S+ was pummeled by \S+
\S+ walked into fire whilst fighting \S+
\S+ was burnt to a crisp whilst fighting \S+
\S+ tried to swim in lava to escape \S+
\S+ drowned whilst trying to escape \S+
\S+ walked into a cactus whilst trying to escape \S+
\S+ was blown up by \S+
`
