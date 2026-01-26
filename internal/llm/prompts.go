package llm

const MemoryEnhancementPrompt string = `The following message was classified as a relevant memory:
%v
Your task is to enhance this message's content to make it more relevant and self-explanatory.
Do **NOT** change the meaning or add content that is not on the message.
Expected Ouput:
Solely the enhanced version of the message, without any confirmation messages.`
