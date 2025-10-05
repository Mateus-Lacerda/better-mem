"""
Locust para testes de escalabilidade do Better-Mem
Simula o comportamento padrão de usuários enviando mensagens e fazendo fetch de memórias.
"""

import random
import time
from locust import HttpUser, task, between


class UltraRemUser(HttpUser):
    wait_time = between(20, 30)
    
    def on_start(self):
        """Executado quando um usuário inicia"""
        # Gera um chat_id único para este usuário
        self.chat_id = f"test_chat_{random.randint(1000, 9999)}"
        self.client.post(
            "/api/v1/chat",
            json={"external_id": self.chat_id},
            name="POST /chat"
        )
        
        # Mensagens triviais (40) - não devem ser consideradas memórias
        self.trivial_messages = [
            "Hi, how are you?",
            "Good morning!",
            "Good afternoon!",
            "Good evening!",
            "Bye!",
            "See you later!",
            "Thank you!",
            "You're welcome!",
            "Please",
            "Excuse me",
            "Sorry",
            "How are things?",
            "What's up?",
            "How's it going?",
            "What time is it?",
            "What day is today?",
            "It's hot today",
            "It's raining",
            "Nice weather",
            "It's cold",
            "Haha",
            "Cool",
            "Interesting",
            "I understand",
            "Right",
            "Ok",
            "Yes",
            "No",
            "Maybe",
            "Could be",
            "Sure",
            "Obviously",
            "True",
            "Exactly",
            "That's right",
            "I agree",
            "I disagree",
            "I don't know",
            "We'll see",
            "We'll figure it out later"
        ]
        
        # Short-term messages (30) - temporary information
        self.short_term_messages = [
            "Today I'm having lunch at the Japanese restaurant",
            "I have a meeting at 3pm",
            "I need to buy milk at the store",
            "I'm going to the doctor tomorrow",
            "I'm watching an interesting movie",
            "I just got home from work",
            "I'm going out for a run now",
            "I'm reading a book about history",
            "I need to finish the report by Friday",
            "I'm traveling this weekend",
            "I have a headache today",
            "I had breakfast late today",
            "Traffic is terrible",
            "I forgot to bring my umbrella",
            "I'm waiting for the bus",
            "I just left the gym",
            "I'm having dinner with my parents today",
            "I need to go to the bank tomorrow",
            "I'm organizing my room",
            "I'm going to the supermarket after work",
            "I'm listening to classical music",
            "I need to study for the exam",
            "I'm going to the movies with friends",
            "I'm cooking pasta",
            "I need to do laundry",
            "I'm going to sleep early today",
            "I'm working from home",
            "I need to call my mom",
            "I'm walking the dog",
            "I'm organizing my schedule"
        ]
        
        # Long-term messages (10) - important and lasting information
        self.long_term_messages = [
            "I'm allergic to peanuts and seafood",
            "My birthday is March 15th",
            "I work as a software engineer for 5 years",
            "I live in São Paulo since 2018",
            "I have two children, an 8-year-old girl and a 5-year-old boy",
            "I've been vegetarian for 3 years for ethical reasons",
            "I graduated in Computer Science from USP in 2019",
            "My favorite color is blue and I love jazz music",
            "I have type 2 diabetes and need to control sugar",
            "I speak Portuguese, English and Spanish fluently"
        ]
        
        # Combina todas as mensagens
        self.all_messages = (
            self.trivial_messages + 
            self.short_term_messages + 
            self.long_term_messages
        )
    
    @task(3)
    def send_message_and_fetch(self, just_fetch=False):
        """
        Comportamento principal: envia uma mensagem e depois faz fetch
        Peso 3 = mais frequente
        """
        # Seleciona uma mensagem aleatória
        message = random.choice(self.all_messages)
        if not just_fetch:
            
            # 1. Envia a mensagem para classificação
            message_payload = {
                "chat_id": self.chat_id,
                "message": message
            }
            
            with self.client.post(
                "/api/v1/message",
                json=message_payload,
                name="POST /message",
                catch_response=True
            ) as response:
                if response.status_code == 202:
                    response.success()
                else:
                    response.failure(f"Status: {response.status_code}")
            
            # Espera um pouco para simular processamento
            time.sleep(random.uniform(0.5, 2.0))
        
        # 2. Faz fetch das memórias usando a mesma mensagem
        fetch_payload = {
            "text": message,
            "limit": 5,
            "vector_search_limit": 20,
            "vector_search_threshold": 0.6,
            "long_term_threshold": 0.8
        }
        
        with self.client.post(
            f"/api/v1/memory/chat/{self.chat_id}/fetch",
            json=fetch_payload,
            name="POST /memory/fetch",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
                # Log do número de memórias retornadas
                memories = response.json()
                print(f"Fetch retornou {len(memories)} memórias para: '{message[:50]}...'")
            else:
                response.failure(f"Status: {response.status_code}")
    
    @task(1)
    def health_check(self):
        """
        Verifica se a API está funcionando
        Peso 1 = menos frequente
        """
        with self.client.get(
            "/api/v1/health",
            name="GET /health",
            catch_response=True
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Status: {response.status_code}")


# Configurações para diferentes cenários de teste
class LightUser(UltraRemUser):
    """Usuário leve - menos requisições"""
    wait_time = between(30, 60)


class HeavyUser(UltraRemUser):
    """Usuário pesado - mais requisições"""
    wait_time = between(10, 20)
    
    @task(5)  # Mais peso para o comportamento principal
    def send_message_and_fetch(self, just_fetch=False):
        super().send_message_and_fetch(just_fetch=just_fetch)


class FetchOnlyUser(UltraRemUser):
    """Usuário que apenas faz fetch"""
    wait_time = between(20, 30)
    
    @task(1)
    def send_message_and_fetch(self, just_fetch=True):
        super().send_message_and_fetch(just_fetch=just_fetch)

if __name__ == "__main__":
    print("Para executar o Locust:")
    print("locust -f locustfile.py --host=http://localhost:5042")
    print("Depois acesse: http://localhost:8089")
